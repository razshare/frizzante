package main

import (
	"embed"
	"github.com/razshare/frizzante/act"
	frizzanteCli "github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/servers"
	"github.com/razshare/frizzante/views"
)

//go:embed .github
//go:embed makefile
//go:embed app/dist
var iefs embed.FS
var port = 8080
var server = make(chan *servers.Server, 1)

func init() {
	// Cli.
	*frizzanteCli.FlagPlatform = "linux/amd64"
	*frizzanteCli.FlagYes = true
	*frizzanteCli.FlagBun = "bun"

	// Server.
	s := servers.New()
	s.Efs = iefs

	s.Routes = append(
		s.Routes,
		routes.Route{Pattern: "GET /TestServerAddRoute", Handler: func(c *connections.Connection) {
			act.SendMessage(c, "hello")
		}},
		routes.Route{Pattern: "GET /TestConnectionSendStatus", Handler: func(c *connections.Connection) {
			act.SendStatus(c, 201)
			act.SendMessage(c, "ok")
		}},
		routes.Route{Pattern: "GET /TestConnectionSendHeader", Handler: func(c *connections.Connection) {
			act.SendHeader(c, "Content-Type", "application/json")
			act.SendMessage(c, "{}")
		}},
		routes.Route{Pattern: "GET /TestRenderServer", Handler: func(c *connections.Connection) {
			act.SendView(c, views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		routes.Route{Pattern: "GET /TestRenderClient", Handler: func(c *connections.Connection) {
			act.SendView(c, views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	)

	go servers.Start(s)

	server <- s
}
