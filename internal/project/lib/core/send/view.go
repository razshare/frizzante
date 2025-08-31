package send

import (
	"main/lib/core/client"
	_view "main/lib/core/view"
	"strings"

	"github.com/razshare/frizzante/stack"
)

// View sends a view.
func View(client *client.Client, view _view.View) {
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
		Json(client, _view.Data(view))
		return
	}

	if client.Config.Render == nil {
		client.Config.ErrorLog.Println("view render function is missing", stack.Trace())
		return
	}

	var html string
	var err error
	if html, err = client.Config.Render(view); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
	}

	if "" == client.Writer.Header().Get("Content-Type") {
		Header(client, "Content-Type", "text/html")
	}

	Message(client, html)
}
