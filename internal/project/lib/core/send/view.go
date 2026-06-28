package send

import (
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

// View sends a view.
func View(http *scopes.Http, view views.View) {
	header := http.Writer.Header()
	if header.Get("Location") != "" {
		return
	}
	if strings.Contains(http.Request.Header.Get("Accept"), "application/json") {
		if header.Get("Cache-Control") == "" {
			Header(http, "Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		}
		if header.Get("Pragma") == "" {
			Header(http, "Pragma", "no-cache")
		}
		if view.Props == nil {
			view.Props = map[string]string{}
		}
		data := views.NewData(view)
		data.Type = http.Request.Header.Get("X-FrizzanteViewType")
		Json(http, data)
		return
	}
	if http.Render == nil {
		logs.Errorf(
			http,
			"send.View: no render function defined\n%s",
			stack.Trace(),
		)
		return
	}
	data := views.NewData(view)
	data.Type = http.Request.Header.Get("X-FrizzanteViewType")
	var html string
	var err error
	if html, err = http.Render(renders.RenderOptions{
		Efs:      http.Efs,
		View:     view,
		Data:     data,
		ErrorLog: http.ErrorLog,
		InfoLog:  http.InfoLog,
	}); err != nil {
		logs.Errorf(
			http,
			"send.View: failed to render view: %v\n%s",
			err,
			stack.Trace(),
		)
	}
	if http.Writer.Header().Get("Content-Type") == "" {
		Header(http, "Content-Type", "text/html")
	}
	Message(http, html)
}
