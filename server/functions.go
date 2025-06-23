package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/guards"
	"github.com/razshare/frizzante/notifiers"
	"github.com/razshare/frizzante/routes"
	"log"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"
)

var server *Server

func init() {
	Initialize()
}

func Initialize() {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	server = &Server{
		Address:         "127.0.0.1:8080",
		SecureAddress:   "127.0.0.1:8383",
		FormMaxMemory:   4096,
		HttpServer:      nil,
		HttpMux:         http.NewServeMux(),
		Connections:     map[string]*net.Conn{},
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		HeaderMaxMemory: 3 * globals.MB,
		Certificate:     "",
		Key:             "",
		Notifier:        notifiers.New(),
		WsUpgrader:      upgrader,
		ViewRoot:        "app",
		PublicRoot:      "app/dist/client",
		ViewServer:      "app/dist/server.js",
		ViewIndex:       "app/dist/client/index.html",
	}
}

// Start starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func Start() {
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
func Stop() {
	shutdownError := server.HttpServer.Shutdown(context.Background())
	if nil != shutdownError {
		log.Fatal(shutdownError)
	}
}

// AddGuard adds a guard.
func AddGuard(val guards.Guard) *Server {
	server.Guards = append(server.Guards, val)
	return server
}

// AddRoute adds a route.
func AddRoute(val routes.Route) *Server {
	server.HttpMux.HandleFunc(val.Pattern, func(writer http.ResponseWriter, request *http.Request) {
		con := &connections.Connection{
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

		for _, tag := range val.Tags {
			for _, guard := range server.Guards {
				if !slices.Contains(guard.Tags, tag) {
					continue
				}
				allowed := false
				guard.Handler(con, func() { allowed = true })
				if !allowed {
					server.Notifier.SendMessage(
						fmt.Sprintf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", val.Pattern, tag, guard.Name),
					)
					return
				}
			}
		}

		val.Handler(con)
	})
	return server
}

// WithAddress sets the server address.
func WithAddress(val string) {
	server.Address = val
}

// WithSecureAddress sets the server secure address.
func WithSecureAddress(val string) {
	server.SecureAddress = val
}

// WithFormMaxMemory sets the server form max memory.
func WithFormMaxMemory(val int64) {
	server.FormMaxMemory = val
}

// WithReadTimeout sets the server read timeout.
func WithReadTimeout(val time.Duration) {
	server.ReadTimeout = val
}

// WithWriteTimeout sets the server write timeout.
func WithWriteTimeout(val time.Duration) {
	server.WriteTimeout = val
}

// WithHeaderMaxMemory sets the server header max memory.
func WithHeaderMaxMemory(val int) {
	server.HeaderMaxMemory = val
}

// WithCertificate sets the server certificate.
func WithCertificate(val string) {
	server.Certificate = val
}

// WithKey sets the server key.
//
// Must match the certificate.
func WithKey(val string) {
	server.Key = val
}

// WithNotifier sets the server notifier.
func WithNotifier(val *notifiers.Notifier) {
	server.Notifier = val
}

// WithEfs sets the server embedded file system.
func WithEfs(efs embed.FS) {
	server.Efs = efs
}

// WithWsUpgrader sets the server web socket upgrader.
func WithWsUpgrader(val *websocket.Upgrader) {
	server.WsUpgrader = val
}

// WithViewRoot sets the server view root.
func WithViewRoot(val string) {
	server.ViewRoot = val
}

// WithPublicRoot sets the server public root.
func WithPublicRoot(val string) {
	server.PublicRoot = val
}

// WithViewServer sets the server view server.
func WithViewServer(val string) {
	server.ViewServer = val
}

// WithViewIndex sets the server view index.
func WithViewIndex(val string) {
	server.ViewIndex = val
}

// Address get the server Address
func Address() string {
	return server.Address
}

// SecureAddress get the server secure address.
func SecureAddress() string {
	return server.SecureAddress
}

// FormMaxMemory get the server form max memory.
func FormMaxMemory() int64 {
	return server.FormMaxMemory
}

// HttpServer get the server http server.
func HttpServer() *http.Server {
	return server.HttpServer
}

// HttpMux get the server http mux.
func HttpMux() *http.ServeMux {
	return server.HttpMux
}

// Connections get the server connections.
func Connections() map[string]*net.Conn {
	return server.Connections
}

// ReadTimeout get the server read timeout.
func ReadTimeout() time.Duration {
	return server.ReadTimeout
}

// WriteTimeout get the server write timeout.
func WriteTimeout() time.Duration {
	return server.WriteTimeout
}

// HeaderMaxMemory get the server header max memory.
func HeaderMaxMemory() int {
	return server.HeaderMaxMemory
}

// Certificate get the server certificate.
func Certificate() string {
	return server.Certificate
}

// Key get the server key.
func Key() string {
	return server.Key
}

// Notifier get the server notifier.
func Notifier() *notifiers.Notifier {
	return server.Notifier
}

// Efs get the server embedded file system.
func Efs() embed.FS {
	return server.Efs
}

// ViewRoot get the server view root.
func ViewRoot() string {
	return server.ViewRoot
}

// PublicRoot get the server public root.
func PublicRoot() string {
	return server.PublicRoot
}

// ViewServer get the server js view server.
func ViewServer() string {
	return server.ViewServer
}

// ViewIndex get the server js view index.
func ViewIndex() string {
	return server.ViewIndex
}

// WsUpgrader get the server web socket upgrader.
func WsUpgrader() *websocket.Upgrader {
	return server.WsUpgrader
}

// Guards get the server guards.
func Guards() []guards.Guard {
	return server.Guards
}
