package server

import (
	"context"
	"errors"
	"main/lib/core/client"
	"net/http"
	"slices"
	"strings"

	"github.com/razshare/frizzante/stack"
)

// Start starts a server from a configuration.
func Start(server *Server) {
	mux := server.Handler.(*http.ServeMux)
	config := &client.Config{
		ErrorLog:   server.ErrorLog,
		InfoLog:    server.InfoLog,
		PublicRoot: server.PublicRoot,
		Efs:        server.Efs,
		Render:     server.Render,
	}
	for _, route := range server.Routes {
		mux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			con := &client.Client{
				Writer:  writer,
				Request: request,
				Config:  config,
				EventId: 1,
				Status:  200,
			}

			for _, tag := range route.Tags {
				for _, guard := range server.Guards {
					if !slices.Contains(guard.Tags, tag) {
						continue
					}
					allow := false
					guard.Handler(con, func() { allow = true })
					if !allow {
						server.InfoLog.Printf("route `%s` tagged with `%d` denied the request because guard `%s` did not pass", route.Pattern, tag, guard.Name)
						return
					}
				}
			}

			route.Handler(con)
		})
	}

	var exit bool

	go func() {
		address := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		server.InfoLog.Printf("server bound to address %s; visit your application at http://%s", server.Addr, address)
		if exit {
			server.InfoLog.Println("cancelling server startup")
			return
		}
		err := http.ListenAndServe(server.Addr, server.Handler)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				server.InfoLog.Println("shutting down server")
				return
			}
			server.ErrorLog.Println(err, stack.Trace())
		}
	}()

	go func() {
		if "" != server.Certificate && "" != server.Key {
			address := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			server.InfoLog.Printf("server bound to address %s; visit your application at https://%s", server.Addr, address)
			if exit {
				server.InfoLog.Println("cancelling server startup")
				return
			}
			err := http.ListenAndServeTLS(server.SecureAddr, server.Certificate, server.Key, server.Handler)
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					server.InfoLog.Println("shutting down server")
					return
				}
				server.ErrorLog.Println(err, stack.Trace())
			}
		}
	}()

	<-server.Channels.Stop
	exit = true

	if err := server.Shutdown(context.Background()); err != nil {
		server.ErrorLog.Println(err)
	}
}
