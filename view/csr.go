package view

import (
	"encoding/json"
	"fmt"
	"github.com/razshare/frizzante/container"
	"os"
	"strings"
)

// RenderClient renders on the client.
func RenderClient(v *View, c *container.Container) (string, error) {
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

	if c.Config.Development {
		var fileNameLocal = strings.ReplaceAll(c.Config.Document, "\\", "/")
		data, rerr := os.ReadFile(fileNameLocal)
		if rerr != nil {
			return "", rerr
		}
		document = string(data)
	} else {
		document = <-c.Channels.Document
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
