package index

import (
	"embed"
	"github.com/razshare/frizzante/view"
	"log"
	"os"
)

// New creates a new index.
func New(efs embed.FS) (c *Index) {
	return &Index{
		Efs:         efs,
		Root:        "app",
		Document:    "app/dist/client/index.html",
		ErrorLog:    log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		Development: os.Getenv("DEV") == "1",
		Channels:    Channels{},
		Render: func(v view.View) (string, error) {
			return Client(c, v)
		},
	}
}
