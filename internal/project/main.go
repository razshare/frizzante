package main

import (
	"embed"
	"main/lib/core/svelte/ssr"
	"main/lib/routes/handlers/fallback"
	"main/lib/routes/handlers/todos"
	"main/lib/routes/handlers/welcome"
	"os"

	"main/lib/core/route"
	"main/lib/core/server"
)

//go:embed app/dist
var efs embed.FS
var srv = server.New()
var dev = os.Getenv("DEV") == "1"
var render = ssr.New(ssr.Config{Efs: efs, Disk: dev})

func main() {
	defer server.Start(srv)
	srv.Efs = efs
	srv.Render = render
	srv.Routes = []route.Route{
		{Pattern: "GET /", Handler: fallback.View},
		{Pattern: "GET /welcome", Handler: welcome.View},
		{Pattern: "GET /todos", Handler: todos.View},
		{Pattern: "GET /check", Handler: todos.Check},
		{Pattern: "GET /uncheck", Handler: todos.Uncheck},
		{Pattern: "GET /add", Handler: todos.Add},
		{Pattern: "GET /remove", Handler: todos.Remove},
	}
}
