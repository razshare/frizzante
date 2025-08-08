package view

import "github.com/razshare/frizzante/container"

// RenderHeadless renders only the body of the view on the server.
func RenderHeadless(v *View, c *container.Container) (string, error) {
	_, body, jsError := container.Render(c, map[string]any{
		"name":       v.Name,
		"data":       v.Data,
		"renderMode": v.RenderMode,
	})
	if jsError != nil {
		return "", jsError
	}
	return body, jsError
}
