package views

import (
	"encoding/json"
	"fmt"
	containers "github.com/razshare/frizzante/apps"
	"github.com/razshare/frizzante/globals"
	"os"
	"strings"
)

// RenderClient renders on the client.
func RenderClient(view *View, app *containers.App, config containers.Config) (string, error) {
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

	if config.Development {
		var fileNameLocal = strings.ReplaceAll(config.Document, "\\", "/")
		data, rerr := os.ReadFile(fileNameLocal)
		if rerr != nil {
			return "", rerr
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
			marshaledProps,
		),
		1,
	), nil
}

// RenderServer renders on the server.
func RenderServer(view *View, app *containers.App, config containers.Config) (string, error) {
	head, body, jsError := containers.ExecuteServerJs(app, config, map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}

	var document string

	if config.Development {
		var fnf = strings.ReplaceAll(config.Document, "\\", "/")
		data, rerr := os.ReadFile(fnf)
		if rerr != nil {
			return "", rerr
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
func RenderHeadless(view *View, app *containers.App, config containers.Config) (string, error) {
	_, body, jsError := containers.ExecuteServerJs(app, config, map[string]any{
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
func RenderFull(view *View, app *containers.App, config containers.Config) (string, error) {
	id := "svelte-app"

	props := map[string]any{
		"name":       view.Name,
		"data":       view.Data,
		"renderMode": view.RenderMode,
	}

	head, body, jsError := containers.ExecuteServerJs(app, config, props)
	if jsError != nil {
		return "", jsError
	}

	var document string

	if config.Development {
		var fnf = strings.ReplaceAll(config.Document, "\\", "/")
		data, rerr := os.ReadFile(fnf)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-app.Document
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
func Render(view *View, app *containers.App, config containers.Config) (string, error) {
	if view.RenderMode == RenderModeFull {
		return RenderFull(view, app, config)
	}

	if view.RenderMode == RenderModeServer {
		return RenderServer(view, app, config)
	}

	if view.RenderMode == RenderModeClient {
		return RenderClient(view, app, config)
	}

	return RenderHeadless(view, app, config)
}
