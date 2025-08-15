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
		if "" == c.Writer.Header().Get("Cache-Control") {
			Header(c, "Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		}
		if "" == c.Writer.Header().Get("Pragma") {
			Header(c, "Pragma", "no-cache")
		}
		if v.Props == nil {
			v.Props = map[string]any{}
		}
		Json(c, view.Data(v))
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
