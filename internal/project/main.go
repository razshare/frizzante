package main

import (
	"embed"
	"errors"
	"log"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/databases"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/servers"
	"github.com/razshare/frizzante/internal/project/lib/core/ssr"
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
var _, queries, databaseError = databases.Connect()
var render = ssr.New(ssr.Options{
	Efs:      efs,
	ErrorLog: errorLog,
	InfoLog:  infoLog,
	Limit:    1,
})
var appRoutes = []routes.Route{
	{Pattern: "GET /", Handler: fallback.View(render, efs)},
	{Pattern: "GET /welcome", Handler: welcome.View(render)},
	{Pattern: "GET /todos", Handler: todos.View(queries, render)},
	{Pattern: "POST /toggle", Handler: todos.Toggle(queries)},
	{Pattern: "POST /add", Handler: todos.Add(queries)},
	{Pattern: "POST /remove", Handler: todos.Remove(queries)},
}
var startError = servers.Start(servers.StartOptions{
	ErrorLog: errorLog,
	InfoLog:  infoLog,
	Routes:   appRoutes,
})

func main() {
	if err := errors.Join(databaseError, startError); err != nil {
		log.Fatal(err)
	}
}
