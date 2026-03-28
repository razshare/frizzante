package servers

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"strings"
	"syscall"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/values"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

// Start starts a server from a configuration.
func Start(server *Server) (err error) {
	var render renders.Render
	if render = server.Render; render == nil {
		err = errors.New("no render function found")
		return
	}
	background := context.Background()
	signalContext, stop := signal.NotifyContext(background, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	go func() {
		<-signalContext.Done()
		server.InfoLog.Println("shutting server down gracefully...")
		if cerr := server.Shutdown(background); cerr != nil {
			if err == nil {
				err = cerr
			}
			return
		}
	}()
	handler := server.Handler.(*http.ServeMux)
	for _, route := range server.Routes {
		handler.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			if errLocal := server.Cors.Check(request); errLocal != nil {
				server.ErrorLog.Printf(
				"servers.Start: CORS check failed: %v",
				errLocal,
			)
				return
			}
			client := &clients.Client{
				Writer:  writer,
				Request: *request,
				Options: clients.Options{
					ErrorLog: server.ErrorLog,
					InfoLog:  server.InfoLog,
					Efs:      server.Efs,
					Render:   render,
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
				client.Channels.End <- values.None
			}
		})
	}
	if server.Certificate != "" && server.Key != "" {
		address := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		server.InfoLog.Printf("server bound to address %s; visit your application at https://%s", server.Addr, address)
		if err = server.ListenAndServeTLS(server.Certificate, server.Key); err != nil {

			if errors.Is(err, http.ErrServerClosed) {
				err = nil
				server.InfoLog.Println("shutting down server")
				return
			}
			return
		}
	} else {
		address := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		server.InfoLog.Printf("server bound to address %s; visit your application at http://%s", server.Addr, address)
		if err = server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				err = nil
				server.InfoLog.Println("shutting down server")
				return
			}
			return
		}
	}
	return
}
