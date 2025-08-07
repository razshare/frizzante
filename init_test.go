package main

import (
	"embed"
	"github.com/razshare/frizzante/act"
	frizzanteCli "github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/servers"
	"github.com/razshare/frizzante/views"
)

//go:embed .github
//go:embed makefile
//go:embed app/dist
var Efs embed.FS
var Port = 8080
var Server = make(chan *servers.Server, 1)

func init() {
	// Cli.
	*frizzanteCli.FlagPlatform = "linux/amd64"
	*frizzanteCli.FlagYes = true
	*frizzanteCli.FlagBun = "bun"

	// Server.
	local := servers.New()
	local.Efs = Efs

	local.Routes = append(
		local.Routes,
		servers.Route{Pattern: "GET /TestServerAddRoute", Handler: func(c *servers.Connection) {
			act.SendMessage(c, "hello")
		}},
		servers.Route{Pattern: "GET /TestConnectionSendStatus", Handler: func(c *servers.Connection) {
			act.SendStatus(c, 201)
			act.SendMessage(c, "ok")
		}},
		servers.Route{Pattern: "GET /TestConnectionSendHeader", Handler: func(c *servers.Connection) {
			act.SendHeader(c, "Content-Type", "application/json")
			act.SendMessage(c, "{}")
		}},
		servers.Route{Pattern: "GET /TestRenderServer", Handler: func(c *servers.Connection) {
			act.SendView(c, views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		servers.Route{Pattern: "GET /TestRenderClient", Handler: func(c *servers.Connection) {
			act.SendView(c, views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	)

	go servers.Start(local)

	Server <- local
}
