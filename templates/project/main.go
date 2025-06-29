package main

import (
	"embed"
	"github.com/razshare/frizzante/web"
	"main/lib/handlers"
)

//go:embed app/dist
var efs embed.FS
var server = web.NewServer()

func main() {
	server.Efs = efs
	server.AddRoute(web.Route{Pattern: "GET /", Handler: handlers.Default})
	server.AddRoute(web.Route{Pattern: "GET /welcome", Handler: handlers.Welcome})
	server.AddRoute(web.Route{Pattern: "GET /todos", Handler: handlers.Todos})
	server.AddRoute(web.Route{Pattern: "GET /check", Handler: handlers.Check})
	server.AddRoute(web.Route{Pattern: "GET /uncheck", Handler: handlers.Uncheck})
	server.AddRoute(web.Route{Pattern: "GET /add", Handler: handlers.Add})
	server.AddRoute(web.Route{Pattern: "GET /remove", Handler: handlers.Remove})
	server.Start()
}
