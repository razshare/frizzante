package frizzante

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

type GuardFunction = func(request *Request, response *Response) bool

type PageConfiguration struct {
	Path         string
	TryFileFirst bool
	Guards       []GuardFunction
}
type PageController interface {
	Configure() PageConfiguration
	Base(request *Request, response *Response)
	Action(request *Request, response *Response)
}

type ApiConfiguration struct {
	Pattern string
	Guards  []GuardFunction
}
type ApiController interface {
	Configure() ApiConfiguration
	Handle(request *Request, response *Response)
}

type Server struct {
	hostName               string
	port                   int
	securePort             int
	multipartFormMaxMemory int64
	server                 *http.Server
	mux                    *http.ServeMux
	connections            map[string]*net.Conn
	readTimeout            time.Duration
	writeTimeout           time.Duration
	maxHeaderBytes         int
	certificate            string
	certificateKey         string
	notifier               *Notifier
	embeddedFileSystem     *embed.FS
	webSocketUpgrader      *websocket.Upgrader
	entryCreated           bool
}

// NewServer creates a server.
func NewServer() *Server {
	notifier := NewNotifier()
	webSocketUpgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	return &Server{
		hostName:               "127.0.0.1",
		port:                   8081,
		securePort:             8383,
		multipartFormMaxMemory: 4096,
		server:                 nil,
		mux:                    http.NewServeMux(),
		connections:            map[string]*net.Conn{},
		readTimeout:            10 * time.Second,
		writeTimeout:           10 * time.Second,
		maxHeaderBytes:         3 * MB,
		certificate:            "",
		certificateKey:         "",
		notifier:               notifier,
		webSocketUpgrader:      webSocketUpgrader,
	}
}

// WithWebSocketReadBufferSize sets the maximum buffer size for each incoming web socket message.
// This will not limit the size of said messages.
func (server *Server) WithWebSocketReadBufferSize(readBufferSize int) {
	server.webSocketUpgrader.ReadBufferSize = readBufferSize
}

// WithWebSocketWriteBufferSize sets the maximum buffer size for each outgoing web socket message.
// This will not limit the size of said messages.
func (server *Server) WithWebSocketWriteBufferSize(writeBufferSize int) {
	server.webSocketUpgrader.WriteBufferSize = writeBufferSize
}

// WithMultipartFormMaxMemory sets the maximum memory for multipart forms before they fall back to disk.
func (server *Server) WithMultipartFormMaxMemory(multipartFormMaxMemory int64) {
	server.multipartFormMaxMemory = multipartFormMaxMemory
}

// WithHostName sets the host name.
func (server *Server) WithHostName(hostName string) {
	server.hostName = hostName
}

// WithPort sets the port.
func (server *Server) WithPort(port int) {
	server.port = port
}

// WithSecurePort sets the secure port.
func (server *Server) WithSecurePort(securePort int) {
	server.securePort = securePort
}

// WithReadTimeout sets the read timeout.
func (server *Server) WithReadTimeout(readTimeout time.Duration) {
	server.readTimeout = readTimeout
}

// WithWriteTimeout sets the write timeout.
func (server *Server) WithWriteTimeout(writeTimeout time.Duration) {
	server.writeTimeout = writeTimeout
}

// WithMaxHeaderBytes sets the maximum allowed bytes in the header of the request.
func (server *Server) WithMaxHeaderBytes(maxHeaderBytes int) {
	server.maxHeaderBytes = maxHeaderBytes
}

// WithCertificateAndKey sets the tls configuration.
func (server *Server) WithCertificateAndKey(certificate string, key string) {
	server.certificate = certificate
	server.certificateKey = key
}

// WithEmbeddedFileSystem sets the embedded file system.
//
// The embedded file system should contain at least directory ".dist" so
// that the server can properly render and serve svelte components.
func (server *Server) WithEmbeddedFileSystem(embeddedFileSystem *embed.FS) {
	server.embeddedFileSystem = embeddedFileSystem
}

// WithNotifier sets the server notifier.
func (server *Server) WithNotifier(notifier *Notifier) {
	server.notifier = notifier
}

// Start starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func (server *Server) Start() {
	logger := log.New(server.notifier.errorFile, "<error>", log.Ltime|log.Llongfile)

	server.server = &http.Server{
		Handler:        server.mux,
		ReadTimeout:    server.readTimeout,
		WriteTimeout:   server.writeTimeout,
		MaxHeaderBytes: server.maxHeaderBytes,
		ErrorLog:       logger,
	}

	var waiter sync.WaitGroup

	waiter.Add(2)

	go func() {
		address := fmt.Sprintf("%s:%d", server.hostName, server.port)
		server.notifier.SendMessage(fmt.Sprintf("listening for requests at http://%s", address))
		serverError := http.ListenAndServe(address, server.mux)
		if nil != serverError {
			if errors.Is(serverError, http.ErrServerClosed) {
				server.notifier.SendMessage("shutting down server")
				return
			}
			log.Fatal(serverError)
		}
	}()

	go func() {
		secureAddress := fmt.Sprintf("%s:%d", server.hostName, server.securePort)
		if "" != server.certificate && "" != server.certificateKey {
			server.notifier.SendMessage(fmt.Sprintf("listening for requests at https://%s", secureAddress))
			serverError := http.ListenAndServeTLS(secureAddress, server.certificate, server.certificateKey, server.mux)
			if nil != serverError {
				if errors.Is(serverError, http.ErrServerClosed) {
					server.notifier.SendMessage("shutting down server")
					return
				}
				log.Fatal(serverError)
			}
		}
	}()

	waiter.Wait()
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

// OnRequest adds request handler.
func (server *Server) OnRequest(pattern string, handler func(request *Request, response *Response)) {
	server.mux.HandleFunc(pattern, func(writer http.ResponseWriter, httpRequest *http.Request) {
		request := &Request{
			server:      server,
			httpRequest: httpRequest,
		}

		httpHeader := writer.Header()

		response := &Response{
			server:                server,
			writer:                &writer,
			lockedStatusAndHeader: false,
			statusCode:            200,
			header:                &httpHeader,
			eventName:             "",
			eventId:               1,
		}

		request.response = response
		response.request = request

		if nil == handler {
			response.SendNotFound()
		}

		handler(request, response)
	})
}

var ids = map[string]string{}

func (server *Server) WithPageController(controller PageController) {
	configuration := controller.Configure()
	reflectedType := reflect.TypeOf(controller)
	id := strings.TrimSuffix(reflectedType.Name(), "Controller")
	ids[id] = configuration.Path
	tryFilesFirst := configuration.TryFileFirst || "/" == configuration.Path
	server.OnRequest("GET "+configuration.Path, func(request *Request, response *Response) {
		if tryFilesFirst {
			response.SendFileOrElse(func() {
				response.id = id
				for _, guard := range configuration.Guards {
					if !guard(request, response) {
						return
					}
				}
				controller.Base(request, response)
			})
		} else {
			response.id = id
			for _, guard := range configuration.Guards {
				if !guard(request, response) {
					return
				}
			}
			controller.Base(request, response)
		}
	})
	server.OnRequest("POST "+configuration.Path, func(request *Request, response *Response) {
		response.id = id
		for _, guard := range configuration.Guards {
			if !guard(request, response) {
				return
			}
		}
		controller.Action(request, response)
	})
}

func (server *Server) WithApiController(controller ApiController) {
	configuration := controller.Configure()
	server.OnRequest(configuration.Pattern, func(request *Request, response *Response) {
		for _, guard := range configuration.Guards {
			if !guard(request, response) {
				return
			}
		}
		controller.Handle(request, response)
	})
}
