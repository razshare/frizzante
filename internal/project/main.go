package main

import (
	"embed"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/route"
	"github.com/razshare/frizzante/internal/project/lib/core/server"
	"github.com/razshare/frizzante/internal/project/lib/core/view/ssr"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/fallback"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/todos"
	"github.com/razshare/frizzante/internal/project/lib/routes/handlers/welcome"
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
