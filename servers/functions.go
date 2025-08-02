package servers

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

func New() *Server {
	return &Server{
		Connections:    map[string]*net.Conn{},
		InfoLog:        log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		SessionArchive: archives.New(filepath.Join(".gen", "sessions")),
		SecureAddr:     "0.0.0.0:8383",
		PublicRoot:     "app/dist/client",
		AppRoot:        "app",
		ServerJs:       "app/dist/server.js",
		IndexHtml:      "app/dist/client/index.html",
		Server: http.Server{
			Addr:           "0.0.0.0:8080",
			Handler:        http.NewServeMux(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 3 * globals.MB,
			ErrorLog:       log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		},
	}
}

// Start starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func (server *Server) Start() {
	mux := server.Handler.(*http.ServeMux)

	for _, route := range server.Routes {
		mux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			con := &connections.Connection{
				Request:        request,
				Writer:         writer,
				Efs:            server.Efs,
				Status:         200,
				EventId:        1,
				PublicRoot:     server.PublicRoot,
				AppRoot:        server.AppRoot,
				ServerJs:       server.ServerJs,
				IndexHtml:      server.IndexHtml,
				ErrorLog:       server.ErrorLog,
				InfoLog:        server.InfoLog,
				SessionArchive: server.SessionArchive,
			}

			for _, tag := range route.Tags {
				for _, guard := range server.Guards {
					if !slices.Contains(guard.Tags, tag) {
						continue
					}
					allowed := false
					guard.Handler(con, func() { allowed = true })
					if !allowed {
						server.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", route.Pattern, tag, guard.Name)
						return
					}
				}
			}

			route.Handler(con)
		})
	}

	var group sync.WaitGroup

	group.Add(2)

	go func() {
		readableAddress := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		server.InfoLog.Printf("listening for requests at http://%s", readableAddress)
		serveError := http.ListenAndServe(server.Addr, server.Handler)
		if serveError != nil {
			if errors.Is(serveError, http.ErrServerClosed) {
				server.InfoLog.Println("shutting down server")
				return
			}
			log.Fatal(serveError)
		}
	}()

	go func() {
		if "" != server.Certificate && "" != server.Key {
			readableAddress := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			server.InfoLog.Printf("listening for requests at https://%s", readableAddress)
			serveError := http.ListenAndServeTLS(server.SecureAddr, server.Certificate, server.Key, server.Handler)
			if serveError != nil {
				if errors.Is(serveError, http.ErrServerClosed) {
					server.InfoLog.Printf("shutting down server")
					return
				}
				log.Fatal(serveError)
			}
		}
	}()

	group.Wait()
}

// Stop attempts to stop the server.
//
// If the shutdown attempt fails, ServerStop crashes the program.
func (server *Server) Stop() {
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
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
