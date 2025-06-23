package libsrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/libcon"
	"github.com/razshare/frizzante/libglobals"
	"github.com/razshare/frizzante/libnotifier"
	"log"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"
)

// NewServer creates a new server.
func NewServer() *Server {
	notifier := libnotifier.NewNotifier()
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Server{
		Address:         "127.0.0.1:8080",
		SecureAddress:   "127.0.0.1:8383",
		FormMaxMemory:   4096,
		HttpServer:      nil,
		HttpMux:         http.NewServeMux(),
		Connections:     map[string]*net.Conn{},
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		HeaderMaxMemory: 3 * libglobals.MB,
		Certificate:     "",
		Key:             "",
		Notifier:        notifier,
		Upgrader:        upgrader,
		ViewRoot:        "app",
		PublicRoot:      "app/dist/client",
		ViewServer:      "app/dist/server.js",
		ViewIndex:       "app/dist/client/index.html",
	}
}

// Start starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func (server *Server) Start() {
	server.HttpServer = &http.Server{
		Handler:        server.HttpMux,
		ReadTimeout:    server.ReadTimeout,
		WriteTimeout:   server.WriteTimeout,
		MaxHeaderBytes: server.HeaderMaxMemory,
		ErrorLog:       server.Notifier.ErrorLogger,
	}

	var group sync.WaitGroup

	group.Add(2)

	go func() {
		server.Notifier.SendMessage(fmt.Sprintf("listening for requests at http://%s", server.Address))
		serverError := http.ListenAndServe(server.Address, server.HttpMux)
		if nil != serverError {
			if errors.Is(serverError, http.ErrServerClosed) {
				server.Notifier.SendMessage("shutting down server")
				return
			}
			log.Fatal(serverError)
		}
	}()

	go func() {
		if "" != server.Certificate && "" != server.Key {
			server.Notifier.SendMessage(fmt.Sprintf("listening for requests at https://%s", server.SecureAddress))
			serverError := http.ListenAndServeTLS(server.SecureAddress, server.Certificate, server.Key, server.HttpMux)
			if nil != serverError {
				if errors.Is(serverError, http.ErrServerClosed) {
					server.Notifier.SendMessage("shutting down server")
					return
				}
				log.Fatal(serverError)
			}
		}
	}()

	group.Wait()
}

// Stop attempts to stop the server.
//
// If the shutdown attempt fails, ServerStop crashes the program.
func (server *Server) Stop() {
	shutdownError := server.HttpServer.Shutdown(context.Background())
	if nil != shutdownError {
		log.Fatal(shutdownError)
	}
}

// AddGuard adds a guard.
func (server *Server) AddGuard(guard Guard) *Server {
	server.Guards = append(server.Guards, guard)
	return server
}

// AddRoute adds a route.
func (server *Server) AddRoute(route Route) *Server {
	server.HttpMux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
		connection := &libcon.Connection{
			PublicRoot: server.PublicRoot,
			Notifier:   server.Notifier,
			Efs:        server.Efs,
			ViewServer: server.ViewServer,
			ViewIndex:  server.ViewIndex,
			ViewRoot:   server.ViewRoot,
			Request:    request,
			Writer:     writer,
			Locked:     false,
			Status:     200,
			Header:     writer.Header(),
			EventName:  "",
			EventId:    1,
		}

		for _, tag := range route.Tags {
			for _, guard := range server.Guards {
				if !slices.Contains(guard.Tags, tag) {
					continue
				}
				allowed := false
				guard.Handler(connection, func() { allowed = true })
				if !allowed {
					server.Notifier.SendMessage(
						fmt.Sprintf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", route.Pattern, tag, guard.Name),
					)
					return
				}
			}
		}

		route.Handler(connection)
	})
	return server
}
