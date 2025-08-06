package servers

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/containers"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/routes"
	"log"
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
		InfoLog:        log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		SessionArchive: archives.NewDiskArchive(filepath.Join(".gen", "sessions")),
		SecureAddr:     "0.0.0.0:8383",
		PublicRoot:     "app/dist/client",
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
	if server.ViewContainer == nil {
		server.ViewContainer = containers.NewViewContainer()
		server.ViewContainer.Efs = server.Efs
	} else {
		if !embeds.IsDirectory(server.ViewContainer.Efs, server.ViewContainer.AppRoot) {
			//trace.
		}
	}

	go server.ViewContainer.Start()

	mux := server.Handler.(*http.ServeMux)

	for _, route := range server.Routes {
		mux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			connection := &connections.Connection{
				Request:        request,
				Writer:         writer,
				Status:         200,
				EventId:        1,
				Efs:            server.Efs,
				PublicRoot:     server.PublicRoot,
				ErrorLog:       server.ErrorLog,
				InfoLog:        server.InfoLog,
				SessionArchive: server.SessionArchive,
				ViewContainer:  server.ViewContainer,
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
