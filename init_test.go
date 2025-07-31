package main

import (
	"embed"
	"fmt"
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

	// Server.
	serverLocal := servers.New()
	serverLocal.Efs = efs
	serverLocal.AddRoute(routes.Route{Pattern: "GET /TestSession", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{Name: "test"}).Start()
		con.SendMessage(fmt.Sprintf("hello %s", session.State.Name))
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "POST /TestSession", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{}).Start()
		defer session.Save()
		session.State.Name = con.ReceiveMessage()
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "GET /TestSessionExpectFail", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{Name: "test"}).Start()
		con.SendMessage(fmt.Sprintf("hello %s", session.State.Name))
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "POST /TestSessionExpectFail", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{}).Start()
		// Without this, session state should not be updated.
		// defer operator.Save(state)
		session.State.Name = con.ReceiveMessage()
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "GET /TestServerAddRoute", Handler: func(con *connections.Connection) {
		con.SendMessage("hello")
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "GET /TestConnectionSendStatus", Handler: func(con *connections.Connection) {
		con.SendStatus(201)
		con.SendMessage("ok")
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "GET /TestConnectionSendHeader", Handler: func(con *connections.Connection) {
		con.SendHeader("Content-Type", "application/json")
		con.SendMessage("{}")
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "GET /TestRenderServer", Handler: func(con *connections.Connection) {
		con.SendView(views.View{
			Name:       "Welcome",
			RenderMode: views.RenderModeServer,
			Data:       map[string]any{"name": "world"},
		})
	}})
	serverLocal.AddRoute(routes.Route{Pattern: "GET /TestRenderClient", Handler: func(con *connections.Connection) {
		con.SendView(views.View{
			Name:       "Welcome",
			RenderMode: views.RenderModeClient,
			Data:       map[string]any{"name": "world"},
		})
	}})

	go serverLocal.Start()

	server <- serverLocal
}
