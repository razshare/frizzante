package server

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/container"
	"github.com/razshare/frizzante/globals"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

// Default creates a new server.
func Default() *Server {
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
		PublicRoot: "app/dist/client",
		Application: container.Configuration{
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
func Start(server *Server) {
	application := container.Start(server.Application, server.Efs)
	defer func() { go func() { application.Stop <- 0 }() }()

	mux := server.Handler.(*http.ServeMux)

	for _, r := range server.Routes {
		mux.HandleFunc(r.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			connection := &Connection{
				EventId:     1,
				Status:      200,
				Writer:      writer,
				Request:     request,
				Application: application,
				Server:      server,
			}

			for _, tag := range r.Tags {
				for _, guard := range server.Guards {
					if !slices.Contains(guard.Tags, tag) {
						continue
					}
					allow := false
					guard.Handler(connection, func() { allow = true })
					if !allow {
						server.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", r.Pattern, tag, guard.Name)
						return
					}
				}
			}

			r.Handler(connection)
		})
	}

	var exit bool

	go func() {
		readableAddress := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		server.InfoLog.Printf("server bound to address %s; visit your application at http://%s", server.Addr, readableAddress)
		if exit {
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
			if exit {
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
	exit = true

	if err := server.Shutdown(context.Background()); err != nil {
		server.ErrorLog.Println(err)
	}
}
