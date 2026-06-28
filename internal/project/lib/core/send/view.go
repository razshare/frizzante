package send

import (
	"net/http"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

// View sends a view.
func View(
	writer http.ResponseWriter,
	request *http.Request,
	render renders.Render,
	view views.View,
) (err error) {
	header := writer.Header()
	if header.Get("Location") != "" {
		return
	}
	if strings.Contains(request.Header.Get("Accept"), "application/json") {
		if header.Get("Cache-Control") == "" {
			header.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		}
		if header.Get("Pragma") == "" {
			header.Set("Pragma", "no-cache")
		}
		if view.Props == nil {
			view.Props = map[string]string{}
		}
		data := views.NewData(view)
		data.Type = request.Header.Get("X-FrizzanteViewType")
		if err = Json(writer, data); err != nil {
			return
		}
		return
	}
	data := views.NewData(view)
	data.Type = request.Header.Get("X-FrizzanteViewType")
	var html string
	if html, err = render(renders.Options{
		View: view,
		Data: data,
	}); err != nil {
		return
	}
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "text/html")
	}
	_, err = writer.Write([]byte(html))
	return
}
