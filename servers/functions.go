package servers

import (
	"context"
	"errors"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
	"sync"
	"time"
)

func New() *Server {
	return &Server{
		Connections:   map[string]*net.Conn{},
		InfoLog:       log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		Address:       "0.0.0.0:8080",
		SecureAddress: "0.0.0.0:8383",
		PublicRoot:    "app/dist/client",
		AppRoot:       "app",
		ServerJs:      "app/dist/server.js",
		IndexHtml:     "app/dist/client/index.html",
		Http: &http.Server{
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
func Start(server *Server) {
	mux := server.Http.Handler.(*http.ServeMux)

	for _, route := range server.Routes {
		mux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			con := &connections.Connection{
				Http:       server.Http,
				Request:    request,
				Writer:     writer,
				Efs:        server.Efs,
				Status:     200,
				EventId:    1,
				PublicRoot: server.PublicRoot,
				AppRoot:    server.AppRoot,
				ServerJs:   server.ServerJs,
				IndexHtml:  server.IndexHtml,
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
		server.InfoLog.Printf("listening for requests at http://%s", server.Address)
		serveError := http.ListenAndServe(server.Address, server.Http.Handler)
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
			server.InfoLog.Printf("listening for requests at https://%s", server.SecureAddress)
			serveError := http.ListenAndServeTLS(server.SecureAddress, server.Certificate, server.Key, server.Http.Handler)
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
func Stop(server *Server) {
	if err := server.Http.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
