package servers

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/channels"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
	"github.com/razshare/frizzante/internal/project/lib/core/views/render"
)

// Start starts a server from a configuration.
func Start(server *Server) {
	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.Handler,
		ReadTimeout:    server.ReadTimeout,
		WriteTimeout:   server.WriteTimeout,
		MaxHeaderBytes: server.MaxHeaderBytes,
		ErrorLog:       server.ErrorLog,
	}

	var err error
	var render_ render.Render
	if render_, err = render.New(); err != nil {
		httpServer.ErrorLog.Println(err)
		return
	}
	handler := server.Handler.(*http.ServeMux)
	for _, route := range server.Routes {
		handler.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			if err = server.Cors.Check(request); err != nil {
				server.ErrorLog.Println(err)
				return
			}
			client := &clients.Client{
				Writer:  writer,
				Request: *request,
				Options: clients.Options{
					ErrorLog: server.ErrorLog,
					InfoLog:  server.InfoLog,
					Efs:      server.Efs,
					Render:   render_,
				},
				EventId: 1,
				Status:  200,
			}
			for _, guard := range route.Guards {
				allow := false
				guard.Handler(client, func() { allow = true })
				if !allow {
					if guard.Name == "" {
						server.InfoLog.Printf("an unnamed guard blocked the request on route %s", route.Pattern)
					} else {
						server.InfoLog.Printf("guard %s blocked the request on route %s", guard.Name, route.Pattern)
					}
					return
				}
			}

			route.Handler(client)

			if client.Channels.End != nil {
				client.Channels.End <- channels.Nothing
			}
		})
	}
	var exit bool
	go func() {
		address := strings.Replace(httpServer.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		server.InfoLog.Printf("server bound to address %s; visit your application at http://%s", httpServer.Addr, address)
		if exit {
			server.InfoLog.Println("cancelling server startup")
			return
		}
		go func() {
			time.Sleep(10 * time.Millisecond)
			if server.Channels.Start != nil {
				server.Channels.Start <- channels.Nothing
			}
		}()
		if err = httpServer.ListenAndServe(); err != nil {
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
			go func() {
				time.Sleep(10 * time.Millisecond)
				if server.Channels.Start != nil {
					server.Channels.Start <- channels.Nothing
				}
			}()
			if err = httpServer.ListenAndServeTLS(server.Certificate, server.Key); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					server.InfoLog.Println("shutting down server")
					return
				}
				httpServer.ErrorLog.Println(err, stack.Trace())
				os.Exit(1)
			}
		}
	}()
	<-server.Channels.End
	exit = true
	if err = httpServer.Shutdown(context.Background()); err != nil {
		httpServer.ErrorLog.Println(err)
	}
}
