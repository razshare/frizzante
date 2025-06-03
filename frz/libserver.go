package frz

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	address         string
	secureAddress   string
	formMaxMemory   int64
	server          *http.Server
	mux             *http.ServeMux
	connections     map[string]*net.Conn
	readTimeout     time.Duration
	writeTimeout    time.Duration
	headerMaxMemory int
	certificate     string
	key             string
	notifier        *Notifier
	dist            embed.FS
	upgrader        *websocket.Upgrader
	guards          []Guard
}

func NewServer() *Server {
	notifier := NewNotifier()
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Server{
		address:         "127.0.0.1:8080",
		secureAddress:   "127.0.0.1:8383",
		formMaxMemory:   4096,
		server:          nil,
		mux:             http.NewServeMux(),
		connections:     map[string]*net.Conn{},
		readTimeout:     10 * time.Second,
		writeTimeout:    10 * time.Second,
		headerMaxMemory: 3 * MB,
		certificate:     "",
		key:             "",
		notifier:        notifier,
		upgrader:        upgrader,
	}
}

// WithWsMaxRedMemory sets the maximum buffer size for each incoming web socket message.
// This will not limit the size of said messages.
func (server *Server) WithWsMaxRedMemory(size int) *Server {
	server.upgrader.ReadBufferSize = size
	return server
}

// WithWsMaxWriteMemory sets the maximum buffer size for each outgoing web socket message.
// This will not limit the size of said messages.
func (server *Server) WithWsMaxWriteMemory(size int) *Server {
	server.upgrader.WriteBufferSize = size
	return server
}

// WithFormMaxMemory sets the maximum memory for multipart forms before they fall back to disk.
func (server *Server) WithFormMaxMemory(size int64) *Server {
	server.formMaxMemory = size
	return server
}

// WithAddress sets the address.
func (server *Server) WithAddress(address string) *Server {
	server.address = address
	return server
}

// WithSecureAddress sets the secure address.
func (server *Server) WithSecureAddress(secureAddress string) *Server {
	server.secureAddress = secureAddress
	return server
}

// WithReadTimeout sets the read timeout.
func (server *Server) WithReadTimeout(timeout time.Duration) *Server {
	server.readTimeout = timeout
	return server
}

// WithWriteTimeout sets the write timeout.
func (server *Server) WithWriteTimeout(timeout time.Duration) *Server {
	server.writeTimeout = timeout
	return server
}

// WithHeaderMaxMemory sets the maximum allowed bytes in the header of the request.
func (server *Server) WithHeaderMaxMemory(maxHeaderBytes int) *Server {
	server.headerMaxMemory = maxHeaderBytes
	return server
}

// WithCertificate sets certificate and key.
func (server *Server) WithCertificate(certificate string, key string) *Server {
	server.certificate = certificate
	server.key = key
	return server
}

// WithNotifier sets the notifier.
func (server *Server) WithNotifier(notifier *Notifier) *Server {
	server.notifier = notifier
	return server
}

// WithDist sets the dist directory.
func (server *Server) WithDist(dist embed.FS) *Server {
	server.dist = dist
	return server
}

// Start starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func (server *Server) Start() {
	server.server = &http.Server{
		Handler:        server.mux,
		ReadTimeout:    server.readTimeout,
		WriteTimeout:   server.writeTimeout,
		MaxHeaderBytes: server.headerMaxMemory,
		ErrorLog:       server.notifier.errorLogger,
	}

	var group sync.WaitGroup

	group.Add(2)

	go func() {
		server.notifier.SendMessage(fmt.Sprintf("listening for requests at http://%s", server.address))
		serverError := http.ListenAndServe(server.address, server.mux)
		if nil != serverError {
			if errors.Is(serverError, http.ErrServerClosed) {
				server.notifier.SendMessage("shutting down server")
				return
			}
			log.Fatal(serverError)
		}
	}()

	go func() {
		if "" != server.certificate && "" != server.key {
			server.notifier.SendMessage(fmt.Sprintf("listening for requests at https://%s", server.secureAddress))
			serverError := http.ListenAndServeTLS(server.secureAddress, server.certificate, server.key, server.mux)
			if nil != serverError {
				if errors.Is(serverError, http.ErrServerClosed) {
					server.notifier.SendMessage("shutting down server")
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
	shutdownError := server.server.Shutdown(context.Background())
	if nil != shutdownError {
		log.Fatal(shutdownError)
	}
}

type Guard struct {
	Name    string
	Handler func(c *Connection, allow func())
	Tags    []string
}
type Route struct {
	Pattern string
	Handler func(c *Connection)
	Tags    []string
}

// AddGuard adds a guard.
func (server *Server) AddGuard(guard Guard) *Server {
	server.guards = append(server.guards, guard)
	return server
}

// AddRoute adds a route.
func (server *Server) AddRoute(route Route) *Server {
	server.mux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
		connection := &Connection{
			server:    server,
			request:   request,
			writer:    writer,
			locked:    false,
			status:    200,
			header:    writer.Header(),
			eventName: "",
			eventId:   1,
		}

		for _, tag := range route.Tags {
			for _, guard := range server.guards {
				if !slices.Contains(guard.Tags, tag) {
					continue
				}
				allowed := false
				guard.Handler(connection, func() { allowed = true })
				if !allowed {
					server.notifier.SendMessage(
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
