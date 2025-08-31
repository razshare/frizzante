package send

import (
	"main/lib/core/client"
	v "main/lib/core/view"
	"strings"

	"github.com/razshare/frizzante/stack"
)

// View sends a view.
func View(client *client.Client, view v.View) {
	if client.Writer.Header().Get("Location") != "" {
		return
	}

	if strings.Contains(client.Request.Header.Get("Accept"), "application/json") {
		if "" == client.Writer.Header().Get("Cache-Control") {
			Header(client, "Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		}
		if "" == client.Writer.Header().Get("Pragma") {
			Header(client, "Pragma", "no-cache")
		}
		if view.Props == nil {
			view.Props = map[string]any{}
		}
		Json(client, v.Data(view))
		return
	}

	if client.Config.Render == nil {
		client.Config.ErrorLog.Println("render function is missing", stack.Trace())
		return
	}

	html, err := client.Config.Render(view)
	if err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
	}

	if "" == client.Writer.Header().Get("Content-Type") {
		Header(client, "Content-Type", "text/html")
	}

	Message(client, html)
}
