package view

import (
	"encoding/json"
	"fmt"
	containers "github.com/razshare/frizzante/container"
	"github.com/razshare/frizzante/globals"
	"os"
	"strings"
)

// RenderClient renders on the client.
func RenderClient(view *View, container *containers.Container, configuration containers.Configuration) (string, error) {
	id := "svelte-app"

	marshaledProps, marshalError := json.Marshal(map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if marshalError != nil {
		return "", marshalError
	}

	var document string

	if configuration.Development {
		var fileNameLocal = strings.ReplaceAll(configuration.Document, "\\", "/")
		data, rerr := os.ReadFile(fileNameLocal)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-container.Document
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
			marshaledProps,
		),
		1,
	), nil
}

// RenderServer renders on the server.
func RenderServer(view *View, container *containers.Container, configuration containers.Configuration) (string, error) {
	head, body, jsError := containers.ExecuteServerJs(container, configuration, map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}

	var document string

	if configuration.Development {
		var fnf = strings.ReplaceAll(configuration.Document, "\\", "/")
		data, rerr := os.ReadFile(fnf)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-container.Document
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
func RenderHeadless(view *View, container *containers.Container, configuration containers.Configuration) (string, error) {
	_, body, jsError := containers.ExecuteServerJs(container, configuration, map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}
	return body, jsError
}

// RenderFull renders on the server and on the client.
func RenderFull(view *View, container *containers.Container, configuration containers.Configuration) (string, error) {
	id := "svelte-app"

	props := map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	}

	head, body, jsError := containers.ExecuteServerJs(container, configuration, props)
	if jsError != nil {
		return "", jsError
	}

	var document string

	if configuration.Development {
		var fnf = strings.ReplaceAll(configuration.Document, "\\", "/")
		data, rerr := os.ReadFile(fnf)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-container.Document
	}

	marshaledProps, marshalError := json.Marshal(props)
	if marshalError != nil {
		return "", marshalError
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
			marshaledProps,
		),
		1,
	), nil
}

// Render renders.
func Render(view *View, container *containers.Container, configuration containers.Configuration) (string, error) {
	if view.RenderMode == RenderModeFull {
		return RenderFull(view, container, configuration)
	}

	if view.RenderMode == RenderModeServer {
		return RenderServer(view, container, configuration)
	}

	if view.RenderMode == RenderModeClient {
		return RenderClient(view, container, configuration)
	}

	return RenderHeadless(view, container, configuration)
}
