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
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Guard = func(req *Request, res *Response) bool

type PageMetadata struct {
	packageName  string
	fileName     string
	functionName string
}

type Controller struct {
	metadata *PageMetadata
	giveWay  bool
	isRoot   bool
	guards   []Guard
	base     func(req *Request, res *Response)
	action   func(req *Request, res *Response)
}

func (page *Controller) findPath() string {
	return "/" + page.findId()
}

func (page *Controller) findId() string {
	parts := strings.SplitN(page.metadata.packageName, strings.Trim(PAGES_ROOT, "/"), 2)
	if len(parts) < 2 {
		log.Fatalf(
			"controllers `%s` must be located under `%s`, but is located under `%s` instead",
			page.metadata.fileName,
			PAGES_ROOT,
			page.metadata.packageName,
		)
	}

	fullId := strings.Trim(parts[1], "/")

	if !strings.HasSuffix(fullId, ".init") {
		log.Fatalf("page `%s` must be created during package initialization, consider creating it inside the `init` function", fullId)
	}

	id := strings.TrimSuffix(fullId, ".init")

	return id
}

func (page *Controller) GiveWay() *Controller {
	page.giveWay = true
	return page
}

func (page *Controller) WithGuard(guard Guard) *Controller {
	page.guards = append(page.guards, guard)
	return page
}

func (page *Controller) WithGuards(guards []Guard) *Controller {
	page.guards = guards
	return page
}

func (page *Controller) WithBase(handler func(req *Request, res *Response)) *Controller {
	page.base = handler
	return page
}

func (page *Controller) WithAction(handler func(req *Request, res *Response)) *Controller {
	page.action = handler
	return page
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

// WithEfs sets the embedded file system.
func (server *Server) WithEfs(dist embed.FS) *Server {
	server.dist = dist
	return server
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
func (server *Server) Start() {
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
func (server *Server) OnRequest(pattern string, guards []Guard, handle func(req *Request, res *Response)) *Server {
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
			response.SendNotFound()
		}

		for _, guard := range guards {
			if !guard(request, response) {
				return
			}
		}

		handle(request, response)
	})
	return server
}

func newControllerMetadata() *PageMetadata {
	pc, file, _, _ := runtime.Caller(2)
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

	return &PageMetadata{
		packageName:  packageName,
		fileName:     fileName,
		functionName: funcName,
	}
}

var ids = map[string]string{}

// LoadController loads the current package as a controller.
func (server *Server) LoadController(configure func(*Controller)) *Server {
	controller := &Controller{
		metadata: newControllerMetadata(),
	}

	if nil != configure {
		configure(controller)
	}

	if nil == controller.base {
		controller.base = func(req *Request, res *Response) {
			res.SendView(NewView(RenderModeFull))
		}
	}

	if nil == controller.action {
		controller.action = func(req *Request, res *Response) {
			res.SendView(NewView(RenderModeFull))
		}
	}

	var isRoot bool
	var controllerPath string
	id := controller.findId()
	if "" == id {
		log.Fatalf(
			"page controller `%s/%s` resolved into a blank id, which is not allowed",
			controller.metadata.packageName, controller.metadata.fileName,
		)
	}

	if strings.ToLower(DEFAULT_PAGE_ID) == strings.ToLower(id) {
		controllerPath = "/"
		isRoot = true
	} else {
		controllerPath = controller.findPath()
		isRoot = "/" == controllerPath
	}

	ids[id] = controllerPath
	tryFilesFirst := controller.giveWay || isRoot
	server.OnRequest("GET "+controllerPath, []Guard{}, func(request *Request, response *Response) {
		if tryFilesFirst {
			response.SendFileOrElse(func() {
				response.id = id
				for _, guard := range controller.guards {
					if !guard(request, response) {
						return
					}
				}
				controller.base(request, response)
			})
		} else {
			response.id = id
			for _, guard := range controller.guards {
				if !guard(request, response) {
					return
				}
			}
			controller.base(request, response)
		}
	})
	server.OnRequest("POST "+controllerPath, []Guard{}, func(request *Request, response *Response) {
		response.id = id
		for _, guard := range controller.guards {
			if !guard(request, response) {
				return
			}
		}
		controller.action(request, response)
	})
	return server
}
