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
)

var IndexHtmlContents string

// ReadIndexHtml reads the contents of the index html document and returns it.
func (view *View) ReadIndexHtml() (string, error) {
	if IndexHtmlContents != "" {
		return IndexHtmlContents, nil
	}

	if files.IsFile(view.Container.IndexHtml) {
		data, readError := os.ReadFile(view.Container.IndexHtml)
		if readError != nil {
			return "", readError
		}

		IndexHtmlContents = string(data)

		return IndexHtmlContents, nil
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

	IndexHtmlContents = string(data)
	return IndexHtmlContents, nil
}

func FixSourceCode(sourceCode string) string {
	return fmt.Sprintf(
		`
			const module={exports:{}}; const render = (function(){
				%s
				return render;
			})()
			render
			`,
		sourceCode,
	)
}

type ViewReadMode int

const (
	ViewReadModeEfs ViewReadMode = 0
	ViewReadModeFs  ViewReadMode = 1
)

type ViewBundleMode int

const (
	ViewBundleModeEsbuild ViewBundleMode = 0
	ViewBundleModeNone    ViewBundleMode = 1
)

// ExecuteServerJs executes the server script.
func (view *View) ExecuteServerJs(properties map[string]any) (string, string, error) {
	if view.Container == nil {
		return "", "", errors.New("view contaienr is nil")
	}

	if view.Container.Mutex == nil {
		return "", "", errors.New("view mutex is nil")
	}

	view.Container.Mutex.Lock()
	defer view.Container.Mutex.Unlock()

	renderPromise, renderError := view.Container.Render(goja.Undefined(), view.Container.Runtime.ToValue(properties))

	if renderError != nil {
		return "", "", renderError
	}

	value := renderPromise.Export().(*goja.Promise).Result().ToObject(view.Container.Runtime)

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

// Contain creates a new Container.
func Contain(configuration *ContainerConfiguration) *Container {
	container := &Container{
		ContainerConfiguration: configuration,
		Runtime:                goja.New(),
		Mutex:                  &sync.Mutex{},
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

	return container
}

// CompileServerJs reads the contents of the server script and returns it.
func (container *Container) CompileServerJs() (*goja.Program, error) {
	var data []byte
	var readError error

	if container.ReadMode == ViewReadModeFs {
		data, readError = os.ReadFile(container.ServerJs)
		if readError != nil {
			return nil, readError
		}
	}

	if container.ReadMode == ViewReadModeEfs {
		data, readError = container.Efs.ReadFile(strings.ReplaceAll(container.ServerJs, "\\", "/"))
		if readError != nil {
			return nil, readError
		}
	}

	var sourceCode string

	if container.BundleMode == ViewBundleModeEsbuild {
		bundledSourceCode, bundleError := javascript.Bundle(container.AppRoot, api.FormatCommonJS, string(data))
		if bundleError != nil {
			return nil, bundleError
		}

		sourceCode = FixSourceCode(bundledSourceCode)
	} else {
		sourceCode = FixSourceCode(string(data))
	}

	program, compileError := goja.Compile("goja", sourceCode, false)
	if compileError != nil {
		return nil, compileError
	}

	return program, nil
}
