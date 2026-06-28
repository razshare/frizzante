package servers

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"strings"
	"syscall"

	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
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
			scope := &scopes.Http{
				Writer:   writer,
				Request:  *request,
				ErrorLog: server.ErrorLog,
				InfoLog:  server.InfoLog,
				Efs:      server.Efs,
				Queries:  server.Queries,
				Render:   render,
				EventId:  1,
				Status:   200,
			}
			for _, guard := range route.Guards {
				allow := false
				guard.Handler(scope, func() { allow = true })
				if !allow {
					if guard.Name == "" {
						server.InfoLog.Printf("an unnamed guard blocked the request on route %s", route.Pattern)
					} else {
						server.InfoLog.Printf("guard %s blocked the request on route %s", guard.Name, route.Pattern)
					}
					return
				}
			}
			route.Handler(scope)
			if scope.WebSocket != nil {
				if cerr := scope.WebSocket.Close(); cerr != nil {
					logs.Errorf(
						scope,
						"send.WsUpgradeWithUpgrader: failed to close WebSocket connection: %v\n%s",
						cerr,
						stack.Trace(),
					)
				}
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
