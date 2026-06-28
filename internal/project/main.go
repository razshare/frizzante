package main

import (
	"embed"
	"log"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/databases"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/servers"
	"github.com/razshare/frizzante/internal/project/lib/core/ssr"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
	"github.com/razshare/frizzante/internal/project/lib/routes/fallback"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
)

//go:generate frizzante clean
//go:generate frizzante configure
//go:embed app/dist
var efs embed.FS
var errorLog = log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
var infoLog = log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
var queries *schema.Queries
var server *servers.Server
var render renders.Render
var err error

func main() {
	if _, queries, err = databases.Connect(); err != nil {
		log.Fatal(err)
	}
	render = ssr.New(1)
	server = servers.New()
	server.InfoLog = infoLog
	server.ErrorLog = errorLog
	server.Routes = []routes.Route{
		{Pattern: "GET /", Handler: fallback.View(render, efs, errorLog, infoLog)},
		{Pattern: "GET /welcome", Handler: welcome.View(render, efs, errorLog, infoLog)},
		{Pattern: "GET /todos", Handler: todos.View(queries, render, efs, errorLog, infoLog)},
		{Pattern: "POST /toggle", Handler: todos.Toggle(queries)},
		{Pattern: "POST /add", Handler: todos.Add(queries)},
		{Pattern: "POST /remove", Handler: todos.Remove(queries)},
	}
	if err = servers.Start(server); err != nil {
		log.Fatal(err)
	}
}
