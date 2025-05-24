package frizzante

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Identifier = func() Identity

type Guard = func(req *Request, res *Response) bool

type PageIdentifier struct {
	value string
}

type PageConfiguration struct {
	Id      Identity
	GiveWay bool
	Guards  []Guard
}

type PageController interface {
	Configure(id Identifier) PageConfiguration
	Base(req *Request, res *Response)
	Action(req *Request, res *Response)
}

type ApiConfiguration struct {
	Id      Identity
	Pattern string
	GiveWay bool
	Guards  []Guard
}

type ApiController interface {
	Configure(meta Identifier) ApiConfiguration
	Handle(req *Request, res *Response)
}

type Identity struct {
	PackageName  string
	FileName     string
	FunctionName string
}

func identify() Identity {
	pc, file, _, _ := runtime.Caller(1)
	_, fileName := path.Split(file)
	descriptor := runtime.FuncForPC(pc)
	parts := strings.Split(descriptor.Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return Identity{
		PackageName:  packageName,
		FileName:     fileName,
		FunctionName: funcName,
	}
}

func (id *Identity) FindPath() string {
	return "/" + id.FindValue()
}

func (id *Identity) FindValue() string {
	if "controller.go" != id.FileName {
		log.Fatalf("controllers must be located inside a file named `controller.go`, received `%s` instead\n", id.FileName)
	}

	parts := strings.SplitN(id.PackageName, strings.Trim(PAGES_ROOT, "/"), 2)
	if len(parts) < 2 {
		log.Fatalf(
			"controllers `%s` must be located under `%s`, but is located under `%s` instead\n",
			id.FileName,
			PAGES_ROOT,
			id.PackageName,
		)
	}

	fullId := strings.Trim(parts[1], "/")
	pageNameParts := strings.SplitN(fullId, ".", 2)

	if len(pageNameParts) < 2 {
		log.Fatalf("page controllers must always be named `Controller`, received empty string in `%s`\n", fullId)
	}

	pageName := pageNameParts[1]

	if "Controller" != pageName && !strings.HasSuffix(pageName, ".Controller") {
		log.Fatalf("page controllers must always be named `Controller`, found `%s` instead in `%s`\n", pageName, fullId)
	}

	value := strings.TrimSuffix(fullId, ".Controller")

	return value
}

type ServerProperties struct {
	Id         string            `json:"id"`
	RenderMode RenderMode        `json:"renderMode"`
	Data       any               `json:"data"`
	Ids        map[string]string `json:"ids"`
}

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
	hasEntry        bool
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

// WithCertificate sets the certificate ands ts key.
func (server *Server) WithCertificate(certificate string, key string) *Server {
	server.certificate = certificate
	server.key = key
	return server
}

// WithNotifier sets the server notifier.
func (server *Server) WithNotifier(notifier *Notifier) *Server {
	server.notifier = notifier
	return server
}

// Start starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func (server *Server) Start(dist embed.FS) {
	server.dist = dist
	logger := log.New(server.notifier.errorFile, "<error>", log.Ltime|log.Llongfile)

	server.server = &http.Server{
		Handler:        server.mux,
		ReadTimeout:    server.readTimeout,
		WriteTimeout:   server.writeTimeout,
		MaxHeaderBytes: server.headerMaxMemory,
		ErrorLog:       logger,
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

// OnRequest adds request handler.
func (server *Server) OnRequest(pattern string, handle func(req *Request, res *Response)) *Server {
	server.mux.HandleFunc(pattern, func(writer http.ResponseWriter, httpRequest *http.Request) {
		request := &Request{
			server:      server,
			httpRequest: httpRequest,
		}

		httpHeader := writer.Header()

		response := &Response{
			server:     server,
			writer:     &writer,
			locked:     false,
			statusCode: 200,
			header:     &httpHeader,
			eventName:  "",
			eventId:    1,
		}

		request.response = response
		response.request = request

		if nil == handle {
			response.SendNotFound("")
		}

		handle(request, response)
	})
	return server
}

var ids = map[string]string{}

func (server *Server) WithPageController(controller PageController) *Server {
	conf := controller.Configure(identify)
	var isRoot bool
	var controllerPath string
	id := conf.Id.FindValue()
	if "" == id {
		log.Fatalf(
			"page `%s/%s` resolved into a blank id, which is not allowed",
			conf.Id.PackageName, conf.Id.FileName,
		)
	}

	if strings.ToLower("any") == strings.ToLower(id) {
		controllerPath = "/"
		isRoot = true
	} else {
		controllerPath = conf.Id.FindPath()
		isRoot = "/" == controllerPath
	}

	ids[id] = controllerPath
	giveWay := conf.GiveWay || isRoot
	server.OnRequest("GET "+controllerPath, func(request *Request, response *Response) {
		if giveWay {
			response.SendFileOrElse(func() {
				response.id = id
				for _, guard := range conf.Guards {
					if !guard(request, response) {
						return
					}
				}
				controller.Base(request, response)
			})
		} else {
			response.id = id
			for _, guard := range conf.Guards {
				if !guard(request, response) {
					return
				}
			}
			controller.Base(request, response)
		}
	})
	server.OnRequest("POST "+controllerPath, func(request *Request, response *Response) {
		response.id = id
		for _, guard := range conf.Guards {
			if !guard(request, response) {
				return
			}
		}
		controller.Action(request, response)
	})
	return server
}

func (server *Server) WithApiController(controller ApiController) *Server {
	conf := controller.Configure(identify)
	conf.Id.FindValue()
	parts := strings.SplitN(conf.Pattern, " ", 2)
	isRoot := len(parts) > 1 && "/" == parts[1]
	giveWay := conf.GiveWay || isRoot
	server.OnRequest(conf.Pattern, func(request *Request, response *Response) {
		if giveWay {
			response.SendFileOrElse(func() {
				for _, guard := range conf.Guards {
					if !guard(request, response) {
						return
					}
				}
				controller.Handle(request, response)
			})
		} else {
			for _, guard := range conf.Guards {
				if !guard(request, response) {
					return
				}
			}
			controller.Handle(request, response)
		}
	})
	return server
}
