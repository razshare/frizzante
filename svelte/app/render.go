package app

import (
	"encoding/json"
	"fmt"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/view"
	"os"
	"strings"
)

// Full renders on the server and on the client.
func Full(c *Config, v view.View) (string, error) {
	id := "svelte-app"

	props := map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	}

	head, body, jsError := RunEntry(c, props)
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
		document = <-c.Channels.Document
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

// Server renders on the server.
func Server(conf *Config, v view.View) (string, error) {
	head, body, jsError := RunEntry(conf, map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}

	var document string

	if conf.Development {
		var fnf = strings.ReplaceAll(conf.Document, "\\", "/")
		data, rerr := os.ReadFile(fnf)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-conf.Channels.Document
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

// Headless renders only the body of the view on the server.
func Headless(conf *Config, v view.View) (string, error) {
	_, body, jsError := RunEntry(conf, map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}
	return body, jsError
}

// Client renders on the client.
func Client(conf *Config, v view.View) (string, error) {
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

	if conf.Development {
		var fileNameLocal = strings.ReplaceAll(conf.Document, "\\", "/")
		data, rerr := os.ReadFile(fileNameLocal)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-conf.Channels.Document
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
