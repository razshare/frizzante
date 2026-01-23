package main

import (
	"embed"
	"log"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/servers"
	"github.com/razshare/frizzante/internal/project/lib/core/ssr"
	"github.com/razshare/frizzante/internal/project/lib/routes/fallback"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
)

//go:generate make clean configure
//go:generate make package
//go:generate make types
//go:embed app/dist
var efs embed.FS

func main() {
	server := servers.New(ssr.NewFunction(1))
	server.Efs = efs
	server.Routes = []routes.Route{
		{Pattern: "GET /", Handler: fallback.View},
		{Pattern: "GET /welcome", Handler: welcome.View},
		{Pattern: "GET /todos", Handler: todos.View},
		{Pattern: "POST /toggle", Handler: todos.Toggle},
		{Pattern: "POST /add", Handler: todos.Add},
		{Pattern: "POST /remove", Handler: todos.Remove},
	}
	if err := servers.Start(server); err != nil {
		log.Fatal(err)
		return
	}
}
