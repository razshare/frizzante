package main

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/route"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/svelte/ssr"
	"github.com/razshare/frizzante/view"
	"os"
)

//go:embed .github
//go:embed makefile
//go:embed app/dist
//go:embed template/project.zip
//go:embed template/lib
var tefs embed.FS
var port = 7878
var serve = make(chan any, 1)

func init() {
	// Server.
	s := server.New()
	s.Addr = fmt.Sprintf("0.0.0.0:%d", port)
	s.Render = ssr.New(ssr.Config{
		Efs:   s.Efs,
		Disk:  os.Getenv("DEV") == "1",
		Limit: 2,
	})
	s.Routes = []route.Route{
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
				Name:   "Welcome",
				Render: view.RenderServer,
				Props:  map[string]any{"name": "world"},
			})
		}},
		{Pattern: "GET /TestRenderClient", Handler: func(c *client.Client) {
			send.View(c, view.View{
				Name:   "Welcome",
				Render: view.RenderClient,
				Props:  map[string]any{"name": "world"},
			})
		}},
	}
	go server.Start(s)
	serve <- 0
}
