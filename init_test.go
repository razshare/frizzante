package main

import (
	"embed"
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
	serverLocal.ViewContainer = views.Contain(&views.ContainerConfiguration{
		Efs:       efs,
		AppRoot:   "app",
		ServerJs:  "app/dist/server.js",
		IndexHtml: "app/dist/client/index.html",
	})
	serverLocal.Routes = append(
		serverLocal.Routes,
		routes.Route{Pattern: "GET /TestSession", Handler: func(connection *connections.Connection) {
			session := sessions.New(connection, State{Name: "test"}).Start()
			connection.SendMessagef("hello %s", session.State.Name)
		}},
		routes.Route{Pattern: "POST /TestSession", Handler: func(connection *connections.Connection) {
			session := sessions.New(connection, State{}).Start()
			defer session.Save()
			session.State.Name = connection.ReceiveMessage()
		}},
		routes.Route{Pattern: "GET /TestSessionExpectFail", Handler: func(connection *connections.Connection) {
			session := sessions.New(connection, State{Name: "test"}).Start()
			connection.SendMessagef("hello %s", session.State.Name)
		}},
		routes.Route{Pattern: "POST /TestSessionExpectFail", Handler: func(connection *connections.Connection) {
			session := sessions.New(connection, State{}).Start().Start()
			// Without this, session state should not be updated.
			// defer operator.Save(state)
			session.State.Name = connection.ReceiveMessage()
		}},
		routes.Route{Pattern: "GET /TestServerAddRoute", Handler: func(connection *connections.Connection) {
			connection.SendMessage("hello")
		}},
		routes.Route{Pattern: "GET /TestConnectionSendStatus", Handler: func(connection *connections.Connection) {
			connection.SendStatus(201)
			connection.SendMessage("ok")
		}},
		routes.Route{Pattern: "GET /TestConnectionSendHeader", Handler: func(connection *connections.Connection) {
			connection.SendHeader("Content-Type", "application/json")
			connection.SendMessage("{}")
		}},
		routes.Route{Pattern: "GET /TestRenderServer", Handler: func(connection *connections.Connection) {
			connection.SendView(views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		routes.Route{Pattern: "GET /TestRenderClient", Handler: func(connection *connections.Connection) {
			connection.SendView(views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	)

	go serverLocal.Start()

	server <- serverLocal
}
