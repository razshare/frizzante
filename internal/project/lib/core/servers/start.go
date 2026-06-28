package servers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Start starts a server.
func Start(options StartOptions) (err error) {
	cors := options.Cors
	routes := options.Routes
	guards := options.Guards
	errorLog := options.ErrorLog
	certificate := options.Certificate
	key := options.Key
	infoLog := options.InfoLog
	if errorLog == nil {
		errorLog = log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime)
	}
	if infoLog == nil {
		infoLog = log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime)
	}
	server := &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        http.NewServeMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2097152, // 2MB,
		ErrorLog:       errorLog,
	}
	background := context.Background()
	sigctx, stop := signal.NotifyContext(background, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	go func() {
		<-sigctx.Done()
		infoLog.Println("shutting server down gracefully...")
		if cerr := server.Shutdown(background); cerr != nil {
			if err == nil {
				err = cerr
			}
			return
		}
	}()
	handler := server.Handler.(*http.ServeMux)
	for _, route := range routes {
		handler.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			if errLocal := cors.Check(request); errLocal != nil {
				server.ErrorLog.Printf(
					"servers.Start: CORS check failed: %v",
					errLocal,
				)
				return
			}
			for _, guard := range guards {
				var allow bool
				if guard.Handler(request, writer, func() { allow = true }); !allow {
					if guard.Name == "" {
						infoLog.Printf("an unnamed guard blocked the request on route %s", route.Pattern)
					} else {
						infoLog.Printf("guard %s blocked the request on route %s", guard.Name, route.Pattern)
					}
					return
				}
			}
			route.Handler(request, writer)
		})
	}
	if certificate != "" && key != "" {
		address := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		infoLog.Printf("server bound to address %s; visit your application at https://%s", server.Addr, address)
		if err = server.ListenAndServeTLS(certificate, key); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				err = nil
				infoLog.Println("shutting down server")
				return
			}
			return
		}
	} else {
		address := strings.Replace(server.Addr, "0.0.0.0:", "127.0.0.1:", 1)
		infoLog.Printf("server bound to address %s; visit your application at http://%s", server.Addr, address)
		if err = server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				err = nil
				infoLog.Println("shutting down server")
				return
			}
			return
		}
	}
	return
}
