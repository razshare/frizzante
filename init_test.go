package main

import (
	"embed"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/route"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/view"
)

//go:embed .github
//go:embed makefile
//go:embed app/dist
var iefs embed.FS
var port = 8080
var ready = make(chan any, 1)

func init() {
	// Cli.
	*cli.FlagPlatform = "linux/amd64"
	*cli.FlagYes = true
	*cli.FlagBun = "bun"

	// Server.
	c := server.Default()
	c.Container.Efs = iefs

	c.Routes = []route.Route{
		{Pattern: "GET /TestRoutes", Handler: func(c *conn.Conn) {
			send.Message(c, "hello")
		}},
		{Pattern: "GET /TestSendStatus", Handler: func(c *conn.Conn) {
			send.Status(c, 201)
			send.Message(c, "ok")
		}},
		{Pattern: "GET /TestSendHeader", Handler: func(c *conn.Conn) {
			send.Header(c, "Content-Type", "application/json")
			send.Message(c, "{}")
		}},
		{Pattern: "GET /TestRenderServer", Handler: func(c *conn.Conn) {
			send.View(c, view.View{
				Name:       "Welcome",
				RenderMode: view.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		{Pattern: "GET /TestRenderClient", Handler: func(c *conn.Conn) {
			send.View(c, view.View{
				Name:       "Welcome",
				RenderMode: view.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	}

	go server.Start(c)

	ready <- 0
}
