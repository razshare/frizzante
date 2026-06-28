package servers

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
)

// Start starts a server from a configuration.
func Start(server *Server, register func(handler *http.ServeMux)) (err error) {
	background := context.Background()
	sigctx, stop := signal.NotifyContext(background, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	go func() {
		<-sigctx.Done()
		server.InfoLog.Println("shutting server down gracefully...")
		if cerr := server.Shutdown(background); cerr != nil {
			if err == nil {
				err = cerr
			}
			return
		}
	}()
	handler := server.Handler.(*http.ServeMux)
	register(handler)
	for _, route := range server.Routes {
		handler.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			if errLocal := server.Cors.Check(request); errLocal != nil {
				server.ErrorLog.Printf(
					"servers.Start: CORS check failed: %v",
					errLocal,
				)
				return
			}
			for _, guard := range route.Guards {
				var allow bool
				if guard.Handler(request, writer, func() { allow = true }); !allow {
					if guard.Name == "" {
						server.InfoLog.Printf("an unnamed guard blocked the request on route %s", route.Pattern)
					} else {
						server.InfoLog.Printf("guard %s blocked the request on route %s", guard.Name, route.Pattern)
					}
					return
				}
			}
			route.Handler(request, writer)
			//if scope.WebSocket != nil {
			//	if cerr := scope.WebSocket.Close(); cerr != nil {
			//		logs.Errorf(
			//			scope,
			//			"send.WsUpgradeWithUpgrader: failed to close WebSocket connection: %v\n%s",
			//			cerr,
			//			stack.Trace(),
			//		)
			//	}
			//}
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
