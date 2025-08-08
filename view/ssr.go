package view

import (
	"fmt"
	"github.com/razshare/frizzante/container"
	"github.com/razshare/frizzante/globals"
	"os"
	"strings"
)

// RenderServer renders on the server.
func RenderServer(v *View, c *container.Container) (string, error) {
	head, body, jsError := container.Render(c, map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	})
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
