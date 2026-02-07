package send

import (
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

// View sends a view.
func View(client *clients.Client, view views.View) {
	header := client.Writer.Header()
	if header.Get("Location") != "" {
		return
	}
	if strings.Contains(client.Request.Header.Get("Accept"), "application/json") {
		if header.Get("Cache-Control") == "" {
			Header(client, "Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		}
		if header.Get("Pragma") == "" {
			Header(client, "Pragma", "no-cache")
		}
		if view.Props == nil {
			view.Props = map[string]string{}
		}
		data := views.NewData(view)
		data.Type = client.Request.Header.Get("X-FrizzanteViewType")
		Json(client, data)
		return
	}
	if client.Options.Render == nil {
		client.Options.ErrorLog.Println("no render function defined", stack.Trace())
		return
	}
	data := views.NewData(view)
	data.Type = client.Request.Header.Get("X-FrizzanteViewType")
	var html string
	var err error
	if html, err = client.Options.Render(renders.RenderOptions{
		Efs:  client.Options.Efs,
		View: view,
		Data: data,
	}); err != nil {
		client.Options.ErrorLog.Println(err, stack.Trace())
	}
	if client.Writer.Header().Get("Content-Type") == "" {
		Header(client, "Content-Type", "text/html")
	}
	Message(client, html)
}
