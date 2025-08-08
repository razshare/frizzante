package view

import "github.com/razshare/frizzante/container"

// Render renders.
func Render(v *View, c *container.Container) (string, error) {
	if v.RenderMode == RenderModeFull {
		return RenderFull(v, c)
	}

	if v.RenderMode == RenderModeServer {
		return RenderServer(v, c)
	}

	if v.RenderMode == RenderModeClient {
		return RenderClient(v, c)
	}

	return RenderHeadless(v, c)
}
