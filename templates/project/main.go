package main

import (
	"embed"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/server"
	"main/lib/handlers"
)

//go:embed app/dist
var efs embed.FS

func main() {
	server.WithEfs(efs)
	server.AddRoute(routes.Route{Pattern: "GET /", Handler: handlers.Default})
	server.AddRoute(routes.Route{Pattern: "GET /welcome", Handler: handlers.Welcome})
	server.AddRoute(routes.Route{Pattern: "GET /todos", Handler: handlers.Todos})
	server.AddRoute(routes.Route{Pattern: "GET /check", Handler: handlers.Check})
	server.AddRoute(routes.Route{Pattern: "GET /uncheck", Handler: handlers.Uncheck})
	server.AddRoute(routes.Route{Pattern: "GET /add", Handler: handlers.Add})
	server.AddRoute(routes.Route{Pattern: "GET /remove", Handler: handlers.Remove})
	server.Start()
}
