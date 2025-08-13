package app

import (
	"embed"
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/view"
	"os"
)

// NewConfig creates a new app configuration.
func NewConfig(efs embed.FS) *Config {
	c := &Config{
		Root:        "app",
		Script:      "app/dist/server.js",
		Document:    "app/dist/client/index.html",
		Parallels:   2,
		Development: os.Getenv("DEV") == "1",
		Channels:    Channels{},
		Server:      server.NewConfig(efs),
	}
	c.Server.Render = NewRender(c)
	return c
}

// NewRender creates a new render function.
func NewRender(c *Config) func(v view.View) (string, error) {
	return func(v view.View) (string, error) {
		if v.RenderMode == view.RenderModeFull {
			return Full(c, v)
		}

		if v.RenderMode == view.RenderModeServer {
			return Server(c, v)
		}

		if v.RenderMode == view.RenderModeClient {
			return Client(c, v)
		}
		return Headless(c, v)
	}
}
