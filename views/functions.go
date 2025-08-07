package views

import (
	"encoding/json"
	"fmt"
	"github.com/razshare/frizzante/apps"
	"github.com/razshare/frizzante/globals"
	"os"
	"strings"
)

// RenderClient renders on the client.
func RenderClient(v *View, a *apps.App, c *apps.Config) (string, error) {
	id := "svelte-app"

	marshaledProps, marshalError := json.Marshal(map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	})
	if marshalError != nil {
		return "", marshalError
	}

	var document string

	if c.Development {
		var fileNameLocal = strings.ReplaceAll(c.Document, "\\", "/")
		data, rerr := os.ReadFile(fileNameLocal)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-a.Document
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
func RenderServer(v *View, a *apps.App, c *apps.Config) (string, error) {
	head, body, jsError := apps.Render(a, c, map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}

	var document string

	if c.Development {
		var fnf = strings.ReplaceAll(c.Document, "\\", "/")
		data, rerr := os.ReadFile(fnf)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-a.Document
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
func RenderHeadless(v *View, a *apps.App, c *apps.Config) (string, error) {
	_, body, jsError := apps.Render(a, c, map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}
	return body, jsError
}

// RenderFull renders on the server and on the client.
func RenderFull(v *View, a *apps.App, c *apps.Config) (string, error) {
	id := "svelte-app"

	props := map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	}

	head, body, jsError := apps.Render(a, c, props)
	if jsError != nil {
		return "", jsError
	}

	var document string

	if c.Development {
		var fnf = strings.ReplaceAll(c.Document, "\\", "/")
		data, rerr := os.ReadFile(fnf)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-a.Document
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
func Render(v *View, a *apps.App, c *apps.Config) (string, error) {
	if v.RenderMode == RenderModeFull {
		return RenderFull(v, a, c)
	}

	if v.RenderMode == RenderModeServer {
		return RenderServer(v, a, c)
	}

	if v.RenderMode == RenderModeClient {
		return RenderClient(v, a, c)
	}

	return RenderHeadless(v, a, c)
}
