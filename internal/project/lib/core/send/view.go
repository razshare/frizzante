package send

import (
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

// View sends a view.
func View(client *clients.Client, view views.View) {
	if client.Writer.Header().Get("Location") != "" {
		return
	}

	if strings.Contains(client.Request.Header.Get("Accept"), "application/json") {
		if client.Writer.Header().Get("Cache-Control") == "" {
			Header(client, "Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		}
		if client.Writer.Header().Get("Pragma") == "" {
			Header(client, "Pragma", "no-cache")
		}
		if view.Props == nil {
			view.Props = map[string]string{}
		}
		Json(client, views.NewData(view))
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

	if client.Writer.Header().Get("Content-Type") == "" {
		Header(client, "Content-Type", "text/html")
	}

	Message(client, html)
}
