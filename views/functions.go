package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/apps"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/javascript"
	"os"
	"strings"
)

// ExecuteServerJs executes the server script.
func (view *View) ExecuteServerJs(app *apps.App, appConfig apps.Configuration, properties map[string]any) (string, string, error) {
	var runtime *goja.Runtime
	var program *goja.Program
	var compileError error

	if appConfig.Development {
		runtime = goja.New()
		var fileNameFixed = strings.ReplaceAll(appConfig.Script, "\\", "/")
		data, readError := os.ReadFile(fileNameFixed)
		if readError != nil {
			return "", "", readError
		}
		bundledSourceCode, bundleError := javascript.Bundle(appConfig.Root, api.FormatCommonJS, string(data))
		if bundleError != nil {
			return "", "", bundleError
		}
		program, compileError = goja.Compile(appConfig.Script, fmt.Sprintf(globals.RenderScriptFormat, bundledSourceCode), false)
		if compileError != nil {
			return "", "", compileError
		}
	} else {
		runtime = <-app.Runtime
		program = <-app.Program
		defer func() { go func() { app.Runtime <- runtime }() }()
		defer func() { go func() { app.Program <- program }() }()
	}

	runResult, runError := runtime.RunProgram(program)
	if runError != nil {
		return "", "", runError
	}

	renderFn, renderIsFn := goja.AssertFunction(runResult)

	if !renderIsFn {
		return "", "", errors.New("render is not a function")
	}

	renderPromise, renderError := renderFn(goja.Undefined(), runtime.ToValue(properties))

	if renderError != nil {
		return "", "", renderError
	}

	value := renderPromise.Export().(*goja.Promise).Result().ToObject(runtime)

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
func (view *View) RenderClient(app *apps.App, appConfig apps.Configuration) (string, error) {
	id := "svelte-app"

	stringifiedProperties, jsoNError := json.Marshal(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsoNError != nil {
		return "", jsoNError
	}

	var document string

	if appConfig.Development {
		var fileNameFixed = strings.ReplaceAll(appConfig.Document, "\\", "/")
		data, readError := os.ReadFile(fileNameFixed)
		if readError != nil {
			return "", readError
		}
		document = string(data)
	} else {
		document = <-app.Document
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					document,
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
func (view *View) RenderServer(app *apps.App, appConfig apps.Configuration) (string, error) {
	head, body, err := view.ExecuteServerJs(app, appConfig, map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if err != nil {
		return "", err
	}

	var document string

	if appConfig.Development {
		var fileNameFixed = strings.ReplaceAll(appConfig.Document, "\\", "/")
		data, readError := os.ReadFile(fileNameFixed)
		if readError != nil {
			return "", readError
		}
		document = string(data)
	} else {
		document = <-app.Document
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					globals.NoScript.ReplaceAllString(document, ""),
					"<!--app-target-->",
					"",
					1,
				),
				"<!--app-body-->",
				fmt.Sprintf("<div id=\"%s\">%s</div>", document, body),
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
func (view *View) RenderHeadless(app *apps.App, appConfig apps.Configuration) (string, error) {
	_, body, err := view.ExecuteServerJs(app, appConfig, map[string]any{
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
func (view *View) RenderFull(app *apps.App, appConfig apps.Configuration) (string, error) {
	id := "svelte-app"

	properties := map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	}

	head, body, err := view.ExecuteServerJs(app, appConfig, properties)
	if err != nil {
		return "", err
	}

	var document string

	if appConfig.Development {
		var fileNameFixed = strings.ReplaceAll(appConfig.Document, "\\", "/")
		data, readError := os.ReadFile(fileNameFixed)
		if readError != nil {
			return "", readError
		}
		document = string(data)
	} else {
		document = <-app.Document
	}

	stringifiedProperties, jsonError := json.Marshal(properties)
	if jsonError != nil {
		return "", jsonError
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					document,
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
func (view *View) Render(app *apps.App, appConfig apps.Configuration) (string, error) {
	if view.RenderMode == RenderModeFull {
		return view.RenderFull(app, appConfig)
	}

	if view.RenderMode == RenderModeServer {
		return view.RenderServer(app, appConfig)
	}

	if view.RenderMode == RenderModeClient {
		return view.RenderClient(app, appConfig)
	}

	return view.RenderHeadless(app, appConfig)
}
