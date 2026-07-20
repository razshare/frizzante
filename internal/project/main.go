package main

import (
	"embed"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/guards"
	"github.com/razshare/frizzante/internal/project/lib/core/negotiate"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/servers"
	"github.com/razshare/frizzante/internal/project/lib/core/ssr"
	"github.com/razshare/frizzante/internal/project/lib/databases"
	"github.com/razshare/frizzante/internal/project/lib/keys"
	"github.com/razshare/frizzante/internal/project/lib/routes/fallback"
	"github.com/razshare/frizzante/internal/project/lib/routes/todos"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
	"github.com/razshare/frizzante/internal/project/lib/schema"
)

//go:generate frizzante clean
//go:generate frizzante configure
//go:embed app/dist
var efs embed.FS
var errorLog = log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
var infoLog = log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
var database, databaseError = databases.Connect()
var queries = schema.New(database)
var render = ssr.New(ssr.Options{
	Efs:      efs,
	ErrorLog: errorLog,
	InfoLog:  infoLog,
	Limit:    1,
})
var session = guards.Guard{
	Name: "session",
	Handler: func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter, allow func()) {
		sessionId, _ := negotiate.SessionId(writer, request)
		session, _ := queries.FindSessionById(request.Context(), sessionId)
		if session.ID == "" {
			if err := queries.AddSessionWithIdAndRoles(request.Context(), schema.AddSessionWithIdAndRolesParams{
				ID:    sessionId,
				Roles: "user",
			}); err != nil {
				errorLog.Printf("error:%v", err)
				return
			}
			session, _ = queries.FindSessionById(request.Context(), sessionId)
		}
		if session.UserID != "guest" {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !strings.Contains(session.Roles, "user") {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		scope[keys.Session] = session
		allow()
	},
}
var appRoutes = []routes.Route{
	{Pattern: "GET /", Handler: fallback.View(efs)},
	{Pattern: "GET /welcome", Handler: welcome.View(render)},
	{Pattern: "GET /todos", Handler: todos.View(queries, render), Guards: []guards.Guard{session}},
	{Pattern: "POST /toggle", Handler: todos.Toggle(queries), Guards: []guards.Guard{session}},
	{Pattern: "POST /add", Handler: todos.Add(queries), Guards: []guards.Guard{session}},
	{Pattern: "POST /remove", Handler: todos.Remove(queries), Guards: []guards.Guard{session}},
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
