package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/javascript"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func RaceForContainer(
	mutex *sync.Mutex,
	containers []*Container,
) (container *Container, unlock func()) {
	for {
		mutex.Lock()
		for _, containerLocal := range containers {
			if containerLocal.Available {
				containerLocal.Available = false
				container = containerLocal
				break
			}
		}

		if container != nil {
			mutex.Unlock()
			break
		}

		mutex.Unlock()

		time.Sleep(time.Millisecond)
	}

	unlock = func() {
		mutex.Lock()
		if container != nil {
			container.Available = true
		}
		mutex.Unlock()
	}

	return
}

// Contain creates a new Container.
func Contain(count int, configuration ContainerConfiguration) []*Container {
	containers := make([]*Container, count)

	for i := 0; i < count; i++ {
		container := &Container{
			ContainerConfiguration: configuration,
			Runtime:                goja.New(),
			Available:              true,
		}

		program, readError := container.CompileServerJs()
		if readError != nil {
			log.Fatal(readError)
		}

		runResult, runError := container.Runtime.RunProgram(program)

		if runError != nil {
			log.Fatal(runError)
		}

		renderFn, renderIsFn := goja.AssertFunction(runResult)

		if !renderIsFn {
			log.Fatal(errors.New("render is not a function"))
		}

		container.Render = renderFn

		containers[i] = container
	}

	return containers
}

// CompileServerJs reads the contents of the server script and returns it.
func (container *Container) CompileServerJs() (*goja.Program, error) {
	var data []byte
	var readError error

	if files.IsFile(container.ServerJs) {
		data, readError = os.ReadFile(container.ServerJs)
		if readError != nil {
			return nil, readError
		}
	} else {
		data, readError = container.Efs.ReadFile(strings.ReplaceAll(container.ServerJs, "\\", "/"))
		if readError != nil {
			return nil, readError
		}
	}

	var sourceCode string

	bundledSourceCode, bundleError := javascript.Bundle(container.AppRoot, api.FormatCommonJS, string(data))
	if bundleError != nil {
		return nil, bundleError
	}

	sourceCode = FixSourceCode(bundledSourceCode)

	program, compileError := goja.Compile("goja", sourceCode, false)
	if compileError != nil {
		return nil, compileError
	}

	return program, nil
}

// ReadIndexHtml reads the contents of the index html document and returns it.
func (view *View) ReadIndexHtml() (string, error) {
	if view.Container.IndexHtmlCache != "" {
		return view.Container.IndexHtmlCache, nil
	}

	if files.IsFile(view.Container.IndexHtml) {
		data, readError := os.ReadFile(view.Container.IndexHtml)
		if readError != nil {
			return "", readError
		}

		view.Container.IndexHtmlCache = string(data)

		return view.Container.IndexHtmlCache, nil
	}

	var data []byte
	fileNameFixed := strings.ReplaceAll(view.Container.IndexHtml, "\\", "/")
	if embeds.IsFile(view.Container.Efs, fileNameFixed) {
		var readError error
		data, readError = view.Container.Efs.ReadFile(fileNameFixed)
		if readError != nil {
			return "", readError
		}
	} else {
		return "", errors.New("view index is missing from the host file system and the embedded file system")
	}

	view.Container.IndexHtmlCache = string(data)
	return view.Container.IndexHtmlCache, nil
}

func FixSourceCode(sourceCode string) string {
	return fmt.Sprintf(
		`
			if (!module) {
				var module={exports:{}}; 
			}

			(function(){
				%s
				return render;
			})()
		`,
		sourceCode,
	)
}

// ExecuteServerJs executes the server script.
func (view *View) ExecuteServerJs(properties map[string]any) (string, string, error) {
	if view.Container == nil {
		return "", "", errors.New("view container is nil")
	}

	runtime := view.Container.Runtime
	render := view.Container.Render

	renderPromise, renderError := render(goja.Undefined(), runtime.ToValue(properties))

	if renderError != nil {
		return "", "", renderError
	}

	value := renderPromise.Export().(*goja.Promise).Result().ToObject(runtime)

	head := value.Get("head")
	body := value.Get("body")

	var headString string
	var bodyString string

	if head != nil {
		headString = head.String()
	}

	if body != nil {
		bodyString = body.String()
	}

	return headString, bodyString, nil

}

// RenderClient renders on the client.
func (view *View) RenderClient() (string, error) {
	id := "svelte-app"

	stringifiedProperties, jsoNError := json.Marshal(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsoNError != nil {
		return "", jsoNError
	}

	indexHtmlData, indexHtmlDataError := view.ReadIndexHtml()
	if indexHtmlDataError != nil {
		return "", indexHtmlDataError
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					indexHtmlData,
					"<!--app-target-->",
					fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", id),
					1,
				),
				"<!--app-body-->",
				fmt.Sprintf("<div id=\"%s\"></div>", id),
				1,
			),
			"<!--app-head-->",
			"",
			1,
		),
		"<!--app-data-->",
		fmt.Sprintf(
			"<script type=\"application/javascript\">function props(){return %s}</script>",
			stringifiedProperties,
		),
		1,
	), nil
}

// RenderServer renders on the server.
func (view *View) RenderServer() (string, error) {
	head, body, err := view.ExecuteServerJs(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if err != nil {
		return "", err
	}

	index, readError := view.ReadIndexHtml()
	if readError != nil {
		return "", readError
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					globals.NoScript.ReplaceAllString(index, ""),
					"<!--app-target-->",
					"",
					1,
				),
				"<!--app-body-->",
				fmt.Sprintf("<div id=\"%s\">%s</div>", index, body),
				1,
			),
			"<!--app-head-->",
			head,
			1,
		),
		"<!--app-data-->",
		"",
		1,
	), nil
}

// RenderHeadless renders only the body of the view on the server.
func (view *View) RenderHeadless() (string, error) {
	_, body, err := view.ExecuteServerJs(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if err != nil {
		return "", err
	}
	return body, err
}

// RenderFull renders on the server and on the client.
func (view *View) RenderFull() (string, error) {
	id := "svelte-app"

	properties := map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	}

	head, body, err := view.ExecuteServerJs(properties)
	if err != nil {
		return "", err
	}

	index, readError := view.ReadIndexHtml()
	if readError != nil {
		return "", readError
	}

	stringifiedProperties, jsonError := json.Marshal(properties)
	if jsonError != nil {
		return "", jsonError
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					index,
					"<!--app-target-->",
					fmt.Sprintf("<script type=\"application/javascript\">function target(){return document.getElementById(\"%s\")}</script>", id),
					1,
				),
				"<!--app-body-->",
				fmt.Sprintf("<div id=\"%s\">%s</div>", id, body),
				1,
			),
			"<!--app-head-->",
			head,
			1,
		),
		"<!--app-data-->",
		fmt.Sprintf(
			"<script type=\"application/javascript\">function props(){return %s}</script>",
			stringifiedProperties,
		),
		1,
	), nil
}

// Render renders.
func (view *View) Render() (string, error) {
	if view.RenderMode == RenderModeFull {
		return view.RenderFull()
	}

	if view.RenderMode == RenderModeServer {
		return view.RenderServer()
	}

	if view.RenderMode == RenderModeClient {
		return view.RenderClient()
	}

	return view.RenderHeadless()
}
