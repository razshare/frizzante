package main

import (
	"embed"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/servers"
	"github.com/razshare/frizzante/internal/project/lib/routes/fallback"
	"github.com/razshare/frizzante/internal/project/lib/routes/login"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
)

//go:generate make clean configure
//go:generate make package
//go:generate make types
//go:embed app/dist
var efs embed.FS
var server = servers.New()

func main() {
	defer servers.Start(server)
	server.Efs = efs
	server.Routes = []routes.Route{
		{Pattern: "GET /", Handler: fallback.View},
		{Pattern: "GET /welcome", Handler: welcome.View},
		{Pattern: "GET /login", Handler: login.View},
		{Pattern: "POST /login", Handler: login.Submit},
		{Pattern: "GET /todos", Handler: todos.View},
		{Pattern: "GET /toggle", Handler: todos.Toggle},
		{Pattern: "GET /add", Handler: todos.Add},
		{Pattern: "GET /remove", Handler: todos.Remove},
	}
}
