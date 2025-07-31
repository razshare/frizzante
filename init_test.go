package main

import (
	"embed"
	"fmt"
	frizzanteCli "github.com/razshare/frizzante/cli"
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
		servers.Route{Pattern: "GET /TestSession", Handler: func(con *servers.Connection) {
			session := sessions.New(con, State{Name: "test"})
			session.Start()
			con.SendMessage(fmt.Sprintf("hello %s", session.State.Name))
		}},
		servers.Route{Pattern: "POST /TestSession", Handler: func(con *servers.Connection) {
			session := sessions.New(con, State{})
			session.Start()
			defer session.Save()
			session.State.Name = con.ReceiveMessage()
		}},
		servers.Route{Pattern: "GET /TestSessionExpectFail", Handler: func(con *servers.Connection) {
			session := sessions.New(con, State{Name: "test"})
			session.Start()
			con.SendMessage(fmt.Sprintf("hello %s", session.State.Name))
		}},
		servers.Route{Pattern: "POST /TestSessionExpectFail", Handler: func(con *servers.Connection) {
			session := sessions.New(con, State{})
			session.Start()
			// Without this, session state should not be updated.
			// defer operator.Save(state)
			session.State.Name = con.ReceiveMessage()
		}},
		servers.Route{Pattern: "GET /TestServerAddRoute", Handler: func(con *servers.Connection) {
			con.SendMessage("hello")
		}},
		servers.Route{Pattern: "GET /TestConnectionSendStatus", Handler: func(con *servers.Connection) {
			con.SendStatus(201)
			con.SendMessage("ok")
		}},
		servers.Route{Pattern: "GET /TestConnectionSendHeader", Handler: func(con *servers.Connection) {
			con.SendHeader("Content-Type", "application/json")
			con.SendMessage("{}")
		}},
		servers.Route{Pattern: "GET /TestRenderServer", Handler: func(con *servers.Connection) {
			con.SendView(views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeServer,
				Data:       map[string]any{"name": "world"},
			})
		}},
		servers.Route{Pattern: "GET /TestRenderClient", Handler: func(con *servers.Connection) {
			con.SendView(views.View{
				Name:       "Welcome",
				RenderMode: views.RenderModeClient,
				Data:       map[string]any{"name": "world"},
			})
		}},
	)

	go serverLocal.Start()

	server <- serverLocal
}
