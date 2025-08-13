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

	if c.Scope.Render == nil {
		c.Scope.ErrorLog.Println("render is missing", stack.Trace())
		return
	}

	html, err := c.Scope.Render(v)
	if err != nil {
		c.Scope.ErrorLog.Println(err, stack.Trace())
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", "text/html")
	}

	Message(c, html)
}
