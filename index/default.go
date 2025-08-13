package index

import (
	"github.com/razshare/frizzante/view"
	"log"
	"os"
)

// Default creates a new index.
func Default() (c *Config) {
	return &Config{
		Root:        "app",
		Document:    "app/dist/client/index.html",
		InfoLog:     log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		ErrorLog:    log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		Development: os.Getenv("DEV") == "1",
		Channels:    Channels{},
		Render: func(v view.View) (string, error) {
			return Client(c, v)
		},
	}
}
