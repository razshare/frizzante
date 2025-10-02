//go:build !dry

package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/view/render"
)

// Start starts a server from a configuration.
func Start(server *Server) {
	handler := server.Handler.(*http.ServeMux)
	config := &client.Config{
		ErrorLog:   server.ErrorLog,
		InfoLog:    server.InfoLog,
		PublicRoot: server.PublicRoot,
		Efs:        server.Efs,
		Render: render.New(render.Config{
			App:      server.App,
			Efs:      server.Efs,
			InfoLog:  server.InfoLog,
			ErrorLog: server.ErrorLog,
		}),
	}
	for _, route := range server.Routes {
		handler.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			if err := server.Cors.Check(request); err != nil {
				server.ErrorLog.Println(err)
				return
			}
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
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				server.InfoLog.Println("shutting down server")
				return
			}
			server.ErrorLog.Println(err, stack.Trace())
			os.Exit(1)
		}
	}()
	go func() {
		if server.Certificate != "" && server.Key != "" {
			address := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
			server.InfoLog.Printf("server bound to address %s; visit your application at https://%s", server.Addr, address)
			if exit {
				server.InfoLog.Println("cancelling server startup")
				return
			}
			if err := server.ListenAndServeTLS(server.Certificate, server.Key); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					server.InfoLog.Println("shutting down server")
					return
				}
				server.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
		}
	}()
	<-server.Channels.Stop
	exit = true
	if err := server.Shutdown(context.Background()); err != nil {
		server.ErrorLog.Println(err)
	}
}
