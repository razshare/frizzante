package main

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/actions"
	frizzanteCli "github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/servers"
	"github.com/razshare/frizzante/sessions"
	"github.com/razshare/frizzante/views"
)

//go:embed .github
//go:embed makefile
//go:embed app/dist
var efs embed.FS
var port = 8080
var server = make(chan *servers.Server, 1)

func init() {
	// Cli.
	*frizzanteCli.FlagPlatform = "linux/amd64"
	*frizzanteCli.FlagYes = true
	*frizzanteCli.FlagBun = "bun"

	// Server.
	serverLocal := servers.New()
	serverLocal.Efs = efs
	serverLocal.Routes = append(
		serverLocal.Routes,
		routes.Route{Pattern: "GET /TestSession", Handler: func(c *connections.Connection) {
			s := sessions.Start(sessions.New(c, State{Name: "test"}))
			actions.SendMessage(c, fmt.Sprintf("hello %s", s.State.Name))
		}},
		routes.Route{Pattern: "POST /TestSession", Handler: func(c *connections.Connection) {
			s := sessions.Start(sessions.New(c, State{}))
			defer sessions.Save(s)
			s.State.Name = actions.ReceiveMessage(c)
		}},
		routes.Route{Pattern: "GET /TestSessionExpectFail", Handler: func(c *connections.Connection) {
			s := sessions.Start(sessions.New(c, State{Name: "test"}))
			actions.SendMessage(c, fmt.Sprintf("hello %s", s.State.Name))
		}},
		routes.Route{Pattern: "POST /TestSessionExpectFail", Handler: func(c *connections.Connection) {
			s := sessions.Start(sessions.New(c, State{}))
			// Without this, session state should not be updated.
			// defer operator.Save(state)
			s.State.Name = actions.ReceiveMessage(c)
		}},
		routes.Route{Pattern: "GET /TestServerAddRoute", Handler: func(c *connections.Connection) {
			actions.SendMessage(c, "hello")
		}},
		routes.Route{Pattern: "GET /TestConnectionSendStatus", Handler: func(c *connections.Connection) {
			actions.SendStatus(c, 201)
			actions.SendMessage(c, "ok")
		}},
		routes.Route{Pattern: "GET /TestConnectionSendHeader", Handler: func(c *connections.Connection) {
			actions.SendHeader(c, "Content-Type", "application/json")
			actions.SendMessage(c, "{}")
		}},
		routes.Route{Pattern: "GET /TestRenderServer", Handler: func(c *connections.Connection) {
			actions.SendView(c, views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		routes.Route{Pattern: "GET /TestRenderClient", Handler: func(c *connections.Connection) {
			actions.SendView(c, views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	)

	go servers.Start(serverLocal)

	server <- serverLocal
}
