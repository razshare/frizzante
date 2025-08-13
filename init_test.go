package main

import (
	"embed"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/route"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/svelte/app"
	"github.com/razshare/frizzante/view"
)

//go:embed .github
//go:embed makefile
//go:embed template/app/dist
var tefs embed.FS
var port = 8080
var ready = make(chan any, 1)

func init() {
	// Cli.
	*state.Platform = "linux/amd64"
	*state.Yes = true
	*state.Bun = "bun"
	*state.App = "template/app"

	// App.
	a := app.NewConfig(tefs)
	a.Document = "template/app/dist/client/index.html"
	a.Script = "template/app/dist/server.js"
	a.Root = "template/app"
	a.Server.PublicRoot = "template/app/dist/client"
	a.Server.Routes = []route.Route{
		{Pattern: "GET /TestRoutes", Handler: func(c *client.Client) {
			send.Message(c, "hello")
		}},
		{Pattern: "GET /TestSendStatus", Handler: func(c *client.Client) {
			send.Status(c, 201)
			send.Message(c, "ok")
		}},
		{Pattern: "GET /TestSendHeader", Handler: func(c *client.Client) {
			send.Header(c, "Content-Type", "application/json")
			send.Message(c, "{}")
		}},
		{Pattern: "GET /TestRenderServer", Handler: func(c *client.Client) {
			send.View(c, view.View{
				Name:       "Welcome",
				RenderMode: view.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		{Pattern: "GET /TestRenderClient", Handler: func(c *client.Client) {
			send.View(c, view.View{
				Name:       "Welcome",
				RenderMode: view.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	}

	go app.Start(a)

	ready <- 0
}
