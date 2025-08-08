package send

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/view"
	"strings"
)

// View sends a view.
func View(c *client.Client, v view.View) {
	if c.Writer.Header().Get("Location") != "" {
		return
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "application/json") {
		if v.Data == nil {
			v.Data = map[string]any{}
		}
		props := map[string]any{
			"name":       v.Name,
			"data":       v.Data,
			"renderMode": v.RenderMode,
		}
		Json(c, props)
		return
	}

	html, err := view.Render(&v, c.Scope.Container)
	if err != nil {
		c.Scope.Container.Config.ErrorLog.Println(err, stack.Trace())
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", "text/html")
	}

	Message(c, html)
}
