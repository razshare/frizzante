package view

import (
	"encoding/json"
	"fmt"
	"github.com/razshare/frizzante/container"
	"os"
	"strings"
)

// RenderFull renders on the server and on the client.
func RenderFull(v *View, c *container.Container) (string, error) {
	id := "svelte-app"

	props := map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	}

	head, body, jsError := container.Render(c, props)
	if jsError != nil {
		return "", jsError
	}

	var document string

	if c.Config.Development {
		var fnf = strings.ReplaceAll(c.Config.Document, "\\", "/")
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
