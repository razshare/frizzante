package container

import (
	"github.com/razshare/frizzante/view"
	"log"
	"os"
)

// Default creates a new container.
func Default() (c *Config) {
	return &Config{
		Parallels:   2,
		Root:        "app",
		Script:      "app/dist/server.js",
		Document:    "app/dist/client/index.html",
		InfoLog:     log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		ErrorLog:    log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		Development: os.Getenv("DEV") == "1",
		Channels:    Channels{},
		Render: func(v view.View) (string, error) {
			if v.RenderMode == view.RenderModeFull {
				return Full(c, v)
			} else if v.RenderMode == view.RenderModeServer {
				return Server(c, v)
			} else if v.RenderMode == view.RenderModeHeadless {
				return Headless(c, v)
			} else if v.RenderMode == view.RenderModeClient {
				return Client(c, v)
			}
			return Server(c, v)
		},
	}
}
