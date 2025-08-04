package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/razshare/frizzante/containers"
	"github.com/razshare/frizzante/globals"
	"strings"
)

// ExecuteServerJs executes the server script.
func (view *View) ExecuteServerJs(container *containers.ViewContainer, properties map[string]any) (string, string, error) {
	Program := <-container.ProgramChannel
	Runtime := <-container.RuntimeChannel

	defer func() { go func() { container.ProgramChannel <- Program }() }()
	defer func() { go func() { container.RuntimeChannel <- Runtime }() }()

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
func (view *View) RenderClient(container *containers.ViewContainer) (string, error) {
	id := "svelte-app"

	stringifiedProperties, jsoNError := json.Marshal(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsoNError != nil {
		return "", jsoNError
	}

	indexHtmlData, indexHtmlDataError := container.ReadIndexHtml()
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
func (view *View) RenderServer(container *containers.ViewContainer) (string, error) {
	head, body, err := view.ExecuteServerJs(container, map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if err != nil {
		return "", err
	}

	index, readError := container.ReadIndexHtml()
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
func (view *View) RenderHeadless(container *containers.ViewContainer) (string, error) {
	_, body, err := view.ExecuteServerJs(container, map[string]any{
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
func (view *View) RenderFull(container *containers.ViewContainer) (string, error) {
	id := "svelte-app"

	properties := map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	}

	head, body, err := view.ExecuteServerJs(container, properties)
	if err != nil {
		return "", err
	}

	index, readError := container.ReadIndexHtml()
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
func (view *View) Render(container *containers.ViewContainer) (string, error) {
	if view.RenderMode == RenderModeFull {
		return view.RenderFull(container)
	}

	if view.RenderMode == RenderModeServer {
		return view.RenderServer(container)
	}

	if view.RenderMode == RenderModeClient {
		return view.RenderClient(container)
	}

	return view.RenderHeadless(container)
}
