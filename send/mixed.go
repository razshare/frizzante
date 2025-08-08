package send

import (
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/view"
	"strings"
)

// View sends a view.
func View(c *conn.Conn, v view.View) {
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

	html, err := view.Render(&v, c.Container)
	if err != nil {
		c.Container.Config.ErrorLog.Println(err, stack.Trace())
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", "text/html")
	}

	Message(c, html)
}
