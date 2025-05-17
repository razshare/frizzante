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
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

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

// ServerCreate creates a server.
func ServerCreate() *Server {
	notifier := NotifierCreate()
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
		entryCreated:           false,
	}
}

// ServerWithWebSocketReadBufferSize sets the maximum buffer size for each incoming web socket message.
// This will not limit the size of said messages.
func ServerWithWebSocketReadBufferSize(self *Server, readBufferSize int) {
	self.webSocketUpgrader.ReadBufferSize = readBufferSize
}

// ServerWithWebSocketWriteBufferSize sets the maximum buffer size for each outgoing web socket message.
// This will not limit the size of said messages.
func ServerWithWebSocketWriteBufferSize(self *Server, writeBufferSize int) {
	self.webSocketUpgrader.WriteBufferSize = writeBufferSize
}

// ServerWithMultipartFormMaxMemory sets the maximum memory for multipart forms before they fall back to disk.
func ServerWithMultipartFormMaxMemory(self *Server, multipartFormMaxMemory int64) {
	self.multipartFormMaxMemory = multipartFormMaxMemory
}

// ServerWithHostName sets the host name.
func ServerWithHostName(self *Server, hostName string) {
	self.hostName = hostName
}

// ServerWithPort sets the port.
func ServerWithPort(self *Server, port int) {
	self.port = port
}

// ServerWithSecurePort sets the secure port.
func ServerWithSecurePort(self *Server, securePort int) {
	self.securePort = securePort
}

// ServerWithReadTimeout sets the read timeout.
func ServerWithReadTimeout(self *Server, readTimeout time.Duration) {
	self.readTimeout = readTimeout
}

// ServerWithWriteTimeout sets the write timeout.
func ServerWithWriteTimeout(self *Server, writeTimeout time.Duration) {
	self.writeTimeout = writeTimeout
}

// ServerWithMaxHeaderBytes sets the maximum allowed bytes in the header of the request.
func ServerWithMaxHeaderBytes(self *Server, maxHeaderBytes int) {
	self.maxHeaderBytes = maxHeaderBytes
}

// ServerWithCertificateAndKey sets the tls configuration.
func ServerWithCertificateAndKey(self *Server, certificate string, key string) {
	self.certificate = certificate
	self.certificateKey = key
}

// ServerWithEmbeddedFileSystem sets the embedded file system.
//
// The embedded file system should contain at least directory ".dist" so
// that the server can properly render and serve svelte components.
func ServerWithEmbeddedFileSystem(self *Server, embeddedFileSystem *embed.FS) {
	self.embeddedFileSystem = embeddedFileSystem
}

// ServerWithNotifier sets the server notifier.
func ServerWithNotifier(self *Server, notifier *Notifier) {
	self.notifier = notifier
}

// ServerStart starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func ServerStart(self *Server) {
	logger := log.New(self.notifier.errorFile, "<error>", log.Ltime|log.Llongfile)

	self.server = &http.Server{
		Handler:        self.mux,
		ReadTimeout:    self.readTimeout,
		WriteTimeout:   self.writeTimeout,
		MaxHeaderBytes: self.maxHeaderBytes,
		ErrorLog:       logger,
	}

	if !self.entryCreated {
		ServerWithApiBuilder(self, func(api *Api) {
			ApiWithPattern(api, "GET /")
			ApiWithRequestHandler(api, func(request *Request, response *Response) {
				ResponseSendStatus(response, 404)
			})
		})
	}

	var waiter sync.WaitGroup

	waiter.Add(2)

	go func() {
		address := fmt.Sprintf("%s:%d", self.hostName, self.port)
		NotifierSendMessage(self.notifier, fmt.Sprintf("listening for requests at http://%s", address))
		err := http.ListenAndServe(address, self.mux)
		if nil != err {
			if errors.Is(err, http.ErrServerClosed) {
				NotifierSendMessage(self.notifier, "shutting down server")
				return
			}
			log.Fatal(err)
		}
	}()

	go func() {
		secureAddress := fmt.Sprintf("%s:%d", self.hostName, self.securePort)
		if "" != self.certificate && "" != self.certificateKey {
			NotifierSendMessage(self.notifier, fmt.Sprintf("listening for requests at https://%s", secureAddress))
			err := http.ListenAndServeTLS(secureAddress, self.certificate, self.certificateKey, self.mux)
			if nil != err {
				if errors.Is(err, http.ErrServerClosed) {
					NotifierSendMessage(self.notifier, "shutting down server")
					return
				}
				log.Fatal(err)
			}
		}
	}()

	waiter.Wait()
}

// ServerStop attempts to stop the server.
//
// If the shutdown attempt fails, ServerStop crashes the program.
func ServerStop(self *Server) {
	err := self.server.Shutdown(context.Background())
	if nil != err {
		log.Fatal(err)
	}
}

var pathParametersPattern = regexp.MustCompile(`{([^{}]+)}`)

// serverMapRoute maps a pattern to a given route.
//
// If the given pattern conflicts with one that is already registered, serverMapRoute crashes the program.
func serverMapRoute(self *Server, pattern string, route *Route) {
	patternParts := strings.Split(pattern, " ")
	patternCounter := len(patternParts)
	isEntry := patternCounter > 1 && strings.HasPrefix(strings.TrimPrefix(filepath.Join(patternParts[1:]...), " "), "/")

	if isEntry && !self.entryCreated {
		self.entryCreated = true
	}

	self.mux.HandleFunc(pattern, func(writer http.ResponseWriter, httpRequest *http.Request) {
		request := &Request{
			server:      self,
			httpRequest: httpRequest,
		}

		httpHeader := writer.Header()

		response := &Response{
			server:                self,
			writer:                &writer,
			lockedStatusAndHeader: false,
			statusCode:            200,
			header:                &httpHeader,
			eventName:             "",
			eventId:               1,
		}

		request.response = response
		response.request = request

		if isEntry {
			ResponseSendEmbeddedFileOrElse(response, func() {
				ResponseSendFileOrElse(response, func() {
					if route.handler != nil {
						if "/favicon.ico" == request.httpRequest.RequestURI {
							ResponseSendNotFound(response)
							return
						}
						route.handler(request, response)

						if !response.lockedStatusAndHeader {
							ResponseSendMessage(response, "")
						}
					}
				})
			})
		} else if route.handler != nil {
			route.handler(request, response)

			if !response.lockedStatusAndHeader {
				ResponseSendMessage(response, "")
			}
		}
	})
}

// ServerWithPageBuilder adds a page.
func ServerWithPageBuilder[T any](self *Server, builder PageBuilder[T]) {
	page := PageCreate[T]()
	PageWithServer(page, self)

	builder(page)

	if "" == page.name {
		NotifierSendError(self.notifier, fmt.Errorf("every page must have a name"))
		return
	}

	if nil == page.view {
		NotifierSendError(self.notifier, fmt.Errorf("every page must have view, page `%s` doesn't", page.name))
		return
	}

	if "" == page.path {
		NotifierSendError(self.notifier, fmt.Errorf("every page must have a path, page `%s` desn't", page.name))
		return
	}

	pages[page.name] = PageMetadata{
		Path:     page.path,
		ViewName: page.view.name,
	}

	serverMapRoute(self, "GET "+page.path, routeCreateWithView(page.base, page.guards, page.view))
	serverMapRoute(self, "POST "+page.path, routeCreateWithView(page.action, page.guards, page.view))
}

// ServerWithApiBuilder adds an api.
func ServerWithApiBuilder(self *Server, builder ApiBuilder) {
	api := ApiCreate()

	builder(api)

	for _, pattern := range api.patterns {
		if "" == pattern {
			NotifierSendError(self.notifier, fmt.Errorf("could not add api because path is empty"))
			return
		}
		serverMapRoute(self, pattern, routeCreate(api.handler, api.guards))
	}
}
