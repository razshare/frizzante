package main

import (
	"embed"
	"github.com/razshare/frizzante/act"
	frizzanteCli "github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/view"
)

//go:embed .github
//go:embed makefile
//go:embed app/dist
var Efs embed.FS
var Port = 8080
var Server = make(chan *server.Server, 1)

func init() {
	// Cli.
	*frizzanteCli.FlagPlatform = "linux/amd64"
	*frizzanteCli.FlagYes = true
	*frizzanteCli.FlagBun = "bun"

	// Server.
	local := server.Default()
	local.Efs = Efs

	local.Routes = append(
		local.Routes,
		server.Route{Pattern: "GET /TestServerAddRoute", Handler: func(c *server.Connection) {
			act.SendMessage(c, "hello")
		}},
		server.Route{Pattern: "GET /TestConnectionSendStatus", Handler: func(c *server.Connection) {
			act.SendStatus(c, 201)
			act.SendMessage(c, "ok")
		}},
		server.Route{Pattern: "GET /TestConnectionSendHeader", Handler: func(c *server.Connection) {
			act.SendHeader(c, "Content-Type", "application/json")
			act.SendMessage(c, "{}")
		}},
		server.Route{Pattern: "GET /TestRenderServer", Handler: func(c *server.Connection) {
			act.SendView(c, view.View{
				Name:       "Welcome",
				RenderMode: view.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		server.Route{Pattern: "GET /TestRenderClient", Handler: func(c *server.Connection) {
			act.SendView(c, view.View{
				Name:       "Welcome",
				RenderMode: view.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	)

	go server.Start(local)

	Server <- local
}
