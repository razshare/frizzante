package views

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/javascript"
	"os"
	"strings"
	"sync"
)

func (view *View) ReadServerJs(efs embed.FS) (string, error) {
	var data []byte
	var readError error

	if files.IsFile(view.ServerJs) {
		data, readError = os.ReadFile(view.ServerJs)
		if readError != nil {
			return "", readError
		}
	} else {
		data, readError = efs.ReadFile(strings.ReplaceAll(view.ServerJs, "\\", "/"))
		if readError != nil {
			return "", readError
		}
	}

	bundledSourceCode, bundleError := javascript.Bundle(view.AppRoot, api.FormatCommonJS, string(data))
	if bundleError != nil {
		return "", bundleError
	}

	return FixSourceCode(bundledSourceCode), nil
}

// ReadIndexHtml reads the contents of the index html document and returns it.
func (view *View) ReadIndexHtml(efs embed.FS) (string, error) {
	if view.IndexHtmlCache != "" {
		return view.IndexHtmlCache, nil
	}

	if files.IsFile(view.IndexHtml) {
		data, readError := os.ReadFile(view.IndexHtml)
		if readError != nil {
			return "", readError
		}

		view.IndexHtmlCache = string(data)

		return view.IndexHtmlCache, nil
	}

	var data []byte
	fileNameFixed := strings.ReplaceAll(view.IndexHtml, "\\", "/")
	if embeds.IsFile(efs, fileNameFixed) {
		var readError error
		data, readError = efs.ReadFile(fileNameFixed)
		if readError != nil {
			return "", readError
		}
	} else {
		return "", errors.New("view index is missing from the host file system and the embedded file system")
	}

	view.IndexHtmlCache = string(data)
	return view.IndexHtmlCache, nil
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

var Program *goja.Program
var Runtime = goja.New()
var Mutex = &sync.Mutex{}

// ExecuteServerJs executes the server script.
func (view *View) ExecuteServerJs(efs embed.FS, properties map[string]any) (string, string, error) {
	Mutex.Lock()
	defer Mutex.Unlock()

	if Runtime == nil {
		Runtime = goja.New()
	}

	if Program == nil {
		sourceCode, readError := view.ReadServerJs(efs)
		if readError != nil {
			return "", "", readError
		}

		program, compileError := goja.Compile("goja", sourceCode, false)
		if compileError != nil {
			return "", "", compileError
		}

		Program = program
	}

	runResult, runError := Runtime.RunProgram(Program)
	if runError != nil {
		return "", "", runError
	}

	renderFn, renderIsFn := goja.AssertFunction(runResult)

	if !renderIsFn {
		return "", "", errors.New("render is not a function")
	}

	renderPromise, renderError := renderFn(goja.Undefined(), Runtime.ToValue(properties))

	if renderError != nil {
		return "", "", renderError
	}

	value := renderPromise.Export().(*goja.Promise).Result().ToObject(Runtime)

	headValue := value.Get("head")
	bodyValue := value.Get("body")

	var headString string
	var bodyString string

	if headValue != nil {
		headString = headValue.String()
	}

	if bodyValue != nil {
		bodyString = bodyValue.String()
	}

	return headString, bodyString, nil
}

// RenderClient renders on the client.
func (view *View) RenderClient(efs embed.FS) (string, error) {
	id := "svelte-app"

	stringifiedProperties, jsoNError := json.Marshal(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsoNError != nil {
		return "", jsoNError
	}

	indexHtmlData, indexHtmlDataError := view.ReadIndexHtml(efs)
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
func (view *View) RenderServer(efs embed.FS) (string, error) {
	head, body, err := view.ExecuteServerJs(efs, map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if err != nil {
		return "", err
	}

	index, readError := view.ReadIndexHtml(efs)
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
func (view *View) RenderHeadless(efs embed.FS) (string, error) {
	_, body, err := view.ExecuteServerJs(efs, map[string]any{
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
func (view *View) RenderFull(efs embed.FS) (string, error) {
	id := "svelte-app"

	properties := map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	}

	head, body, err := view.ExecuteServerJs(efs, properties)
	if err != nil {
		return "", err
	}

	index, readError := view.ReadIndexHtml(efs)
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
func (view *View) Render(efs embed.FS) (string, error) {
	if view.RenderMode == RenderModeFull {
		return view.RenderFull(efs)
	}

	if view.RenderMode == RenderModeServer {
		return view.RenderServer(efs)
	}

	if view.RenderMode == RenderModeClient {
		return view.RenderClient(efs)
	}

	return view.RenderHeadless(efs)
}
