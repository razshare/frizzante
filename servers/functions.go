package servers

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/apps"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

func New() *Server {
	infoLog := log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	return &Server{
		InfoLog:    infoLog,
		SecureAddr: "0.0.0.0:8383",
		Stop:       make(chan any),
		Server: http.Server{
			Addr:           "0.0.0.0:8080",
			Handler:        http.NewServeMux(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 3 * globals.MB,
			ErrorLog:       errorLog,
		},
		SessionArchive: archives.NewDiskArchive(filepath.Join(".gen", "sessions")),
		PublicRoot:     "app/dist/client",
		AppConfiguration: apps.Configuration{
			Root:        "app",
			Script:      "app/dist/server.js",
			Document:    "app/dist/client/index.html",
			Parallels:   2,
			InfoLog:     infoLog,
			ErrorLog:    errorLog,
			Development: os.Getenv("DEV") == "1",
		},
	}
}

// Start starts the server.
func (server *Server) Start() {
	var app *apps.App

	app = apps.Start(server.AppConfiguration, server.Efs)

	defer func() { go func() { app.Stop <- 0 }() }()

	mux := server.Handler.(*http.ServeMux)

	for _, route := range server.Routes {
		mux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			connection := &connections.Connection{
				EventId:          1,
				Status:           200,
				Writer:           writer,
				Request:          request,
				App:              app,
				Efs:              server.Efs,
				InfoLog:          server.InfoLog,
				ErrorLog:         server.ErrorLog,
				PublicRoot:       server.PublicRoot,
				SessionArchive:   server.SessionArchive,
				AppConfiguration: server.AppConfiguration,
			}

			for _, tag := range route.Tags {
				for _, guard := range server.Guards {
					if !slices.Contains(guard.Tags, tag) {
						continue
					}
					allowed := false
					guard.Handler(connection, func() { allowed = true })
					if !allowed {
						server.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", route.Pattern, tag, guard.Name)
						return
					}
				}
			}

			route.Handler(connection)
		})
	}

	var cancelled bool

	go func() {
		readableAddress := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		server.InfoLog.Printf("server bound to address %s; visit your application at http://%s", server.Addr, readableAddress)
		if cancelled {
			server.InfoLog.Printf("cancelling server startup")
			return
		}
		serveError := http.ListenAndServe(server.Addr, server.Handler)
		if serveError != nil {
			if errors.Is(serveError, http.ErrServerClosed) {
				server.InfoLog.Println("shutting down server")
				return
			}
			server.ErrorLog.Println(serveError)
		}
	}()

	go func() {
		if "" != server.Certificate && "" != server.Key {
			readableAddress := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			server.InfoLog.Printf("server bound to address %s; visit your application at https://%s", server.Addr, readableAddress)
			if cancelled {
				server.InfoLog.Printf("cancelling server startup")
				return
			}
			serveError := http.ListenAndServeTLS(server.SecureAddr, server.Certificate, server.Key, server.Handler)
			if serveError != nil {
				if errors.Is(serveError, http.ErrServerClosed) {
					server.InfoLog.Printf("shutting down server")
					return
				}
				server.ErrorLog.Println(serveError)
			}
		}
	}()

	<-server.Stop
	println("cancelled")
	cancelled = true

	if err := server.Shutdown(context.Background()); err != nil {
		server.ErrorLog.Println(err)
	}
}

// AddGuard adds a guard.
//
// Deprecated: append directly to Guards instead.
func (server *Server) AddGuard(guard guards.Guard) *Server {
	server.Guards = append(server.Guards, guard)
	return server
}

// AddRoute adds a guard.
//
// Deprecated: append directly to Routes instead.
func (server *Server) AddRoute(route routes.Route) *Server {
	server.Routes = append(server.Routes, route)
	return server
}
