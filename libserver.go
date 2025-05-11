package frizzante

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
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
	sessions               map[string]*net.Conn
	readTimeout            time.Duration
	writeTimeout           time.Duration
	maxHeaderBytes         int
	certificate            string
	certificateKey         string
	notifier               *Notifier
	embeddedFileSystem     embed.FS
	webSocketUpgrader      *websocket.Upgrader
	entryCreated           bool
	sessionBuilder         any
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
		sessions:               map[string]*net.Conn{},
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

// ServerWithSessionBuilder sets the session builder.
func ServerWithSessionBuilder[T any](self *Server, builder SessionBuilder[T]) {
	self.sessionBuilder = builder
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
func ServerWithEmbeddedFileSystem(self *Server, embeddedFileSystem embed.FS) {
	self.embeddedFileSystem = embeddedFileSystem
}

// ServerWithNotifier sets the server notifier.
func ServerWithNotifier(self *Server, notifier *Notifier) {
	self.notifier = notifier
}

// RequestReceiveCancellation returns a channel that's closed when the request is cancelled.
func RequestReceiveCancellation(self *Request) <-chan struct{} {
	return self.httpRequest.Context().Done()
}

// RequestReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func RequestReceiveCookie(self *Request, key string) string {
	cookie, cookieError := self.httpRequest.Cookie(key)
	if cookieError != nil {
		NotifierSendError(self.server.notifier, cookieError)
		return ""
	}
	value, unescapeError := url.QueryUnescape(cookie.Value)
	if unescapeError != nil {
		return ""
	}

	return value
}

// RequestReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func RequestReceiveMessage(self *Request) string {
	if self.webSocketConn != nil {
		_, readBytes, readError := self.webSocketConn.ReadMessage()
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return ""
		}
		return string(readBytes)
	}

	readBytes, readAllError := io.ReadAll(self.httpRequest.Body)
	if readAllError != nil {
		NotifierSendError(self.server.notifier, readAllError)
		return ""
	}
	return string(readBytes)
}

// RequestReceiveJson reads the message as json and returns the value and a boolean,
// which indicates success or failure.
//
// Compatible with web sockets.
func RequestReceiveJson[T any](self *Request) *T {
	var value T
	if self.webSocketConn != nil {
		jsonError := self.webSocketConn.ReadJSON(value)
		if jsonError != nil {
			NotifierSendError(self.server.notifier, jsonError)
			return nil
		}
		return &value
	}

	readBytes, readAllError := io.ReadAll(self.httpRequest.Body)
	if readAllError != nil {
		NotifierSendError(self.server.notifier, readAllError)
		return nil
	}
	unmarshalError := json.Unmarshal(readBytes, &value)
	if unmarshalError != nil {
		NotifierSendError(self.server.notifier, unmarshalError)
		return nil
	}
	return &value
}

// RequestReceiveForm reads the message as a form and returns the value.
func RequestReceiveForm(self *Request) *url.Values {
	if self.webSocketConn != nil {
		NotifierSendError(self.server.notifier, errors.New("web socket connections cannot receive form payloads"))
		return &url.Values{}
	}

	parseMultipartFormError := self.httpRequest.ParseMultipartForm(self.server.multipartFormMaxMemory)
	if parseMultipartFormError != nil {
		if !errors.Is(parseMultipartFormError, http.ErrNotMultipart) {
			NotifierSendError(self.server.notifier, parseMultipartFormError)
		}

		parseFormError := self.httpRequest.ParseForm()
		if parseFormError != nil {
			NotifierSendError(self.server.notifier, parseFormError)
		}
	}

	return &self.httpRequest.Form
}

// RequestReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func RequestReceiveQuery(self *Request, name string) string {
	return self.httpRequest.URL.Query().Get(name)
}

// RequestReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func RequestReceivePath(self *Request, name string) string {
	return self.httpRequest.PathValue(name)
}

// RequestReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func RequestReceiveHeader(self *Request, key string) string {
	return self.httpRequest.Header.Get(key)
}

// RequestReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func RequestReceiveContentType(self *Request) string {
	return self.httpRequest.Header.Get("Content-Type")
}

func notFoundApi(api *Api) {
	ApiWithPattern(api, "GET /")
	ApiWithRequestHandler(api, func(request *Request, response *Response) {
		ResponseSendStatus(response, 404)
	})
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
		ServerWithApiBuilder(self, notFoundApi)
	}

	var waiter sync.WaitGroup

	waiter.Add(2)

	go func() {
		address := fmt.Sprintf("%s:%d", self.hostName, self.port)
		NotifierSendMessage(self.notifier, fmt.Sprintf("listening for requests at http://%s", address))
		err := http.ListenAndServe(address, self.mux)
		if err != nil {
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
			if err != nil {
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
	if err != nil {
		log.Fatal(err)
	}
}

var pathParametersPattern = regexp.MustCompile(`{([^{}]+)}`)

type Route struct {
	server  *Server
	view    string
	handler func(request *Request, response *Response)
	guards  []func(request *Request, response *Response, pass func())
	mount   func(pattern string)
}

// routeCreate creates a route configuration from a callback function.
func routeCreate(
	handler func(request *Request, response *Response),
	guards []func(request *Request, response *Response, pass func()),
) *Route {
	return &Route{
		view: "",
		handler: func(request *Request, response *Response) {
			pass := 0 == len(guards)
			for _, guard := range guards {
				guard(request, response, func() { pass = true })

				if !pass {
					break
				}
			}

			if pass {
				handler(request, response)
			}
		},
		mount: func(pattern string) {},
	}
}

// routeCreateWithPage creates a route configuration from a callback function, just like routeCreate.
//
// Unlike routeCreate, routeCreateWithView also creates a View, which is used to automatically
// to serve a svelte view after invoking callback.
//
// Generally speaking, you should never manually invoke ResponseSendMessage or similar functions.
//
// However, it is safe to invoke receive functions, like RequestReceiveHeader, RequestReceiveCookie, etc.
func routeCreateWithView(
	handler func(request *Request, response *Response, view *View),
	guards []func(request *Request, response *Response, pass func()),
	view string,
) *Route {
	var pattern string

	return &Route{
		view: view,
		handler: func(
			request *Request,
			response *Response,
		) {
			viewLocal := &View{
				render:     RenderFull,
				data:       map[string]any{},
				name:       view,
				parameters: map[string]string{},
			}

			pass := 0 == len(guards)
			for _, guard := range guards {
				guard(request, response, func() { pass = true })

				if !pass {
					break
				}
			}

			if pass {
				handler(request, response, viewLocal)
			}

			if nil != response.navigate {
				ResponseSendRedirect(response, response.navigate.Location, http.StatusFound)
				return
			}

			if "" != response.header.Get("Location") {
				return
			}

			if nil == viewLocal {
				NotifierSendError(request.server.notifier, fmt.Errorf("svelte page handler `%s` returned a nil page", pattern))
				return
			}

			if nil == viewLocal.data {
				viewLocal.data = map[string]any{}
			}

			if RequestVerifyAccept(request, "application/json") {
				data, marshalError := json.Marshal(viewLocal.data)
				if marshalError != nil {
					NotifierSendError(request.server.notifier, marshalError)
					return
				}
				ResponseSendHeader(response, "Content-Type", "application/json")
				ResponseSendMessage(response, string(data))
				return
			}

			if nil == viewLocal.parameters {
				viewLocal.parameters = map[string]string{}
			}

			for _, name := range pathParametersPattern.FindAllStringSubmatch(pattern, -1) {
				if len(name) < 1 {
					continue
				}
				viewLocal.parameters[name[1]] = request.httpRequest.PathValue(name[1])
			}

			ResponseSendView(response, viewLocal)
		},
		mount: func(patternLocal string) {
			pattern = patternLocal
			patternParts := strings.Split(patternLocal, " ")
			patternCounter := len(patternParts)
			if patternCounter > 1 {
				components[view] = path.Join(patternParts[1:]...)
			}
		},
	}
}

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

	if route.mount != nil {
		route.mount(pattern)
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

		for _, after := range response.after {
			after()
		}
	})
}

type Request struct {
	server        *Server
	response      *Response
	httpRequest   *http.Request
	webSocketConn *websocket.Conn
}

type navigate struct {
	Page       string
	Parameters map[string]string
	Location   string
}

type Response struct {
	server                *Server
	request               *Request
	writer                *http.ResponseWriter
	lockedStatusAndHeader bool
	statusCode            int
	header                *http.Header
	webSocket             *websocket.Conn
	eventName             string
	navigate              *navigate
	eventId               int64
	after                 []func()
	context               map[string]any
}

var pathFieldRegex = regexp.MustCompile(`\{(.*?)}`)

// ResponseSendNavigateWithParameters sends the client an instruction to navigate.
func ResponseSendNavigateWithParameters(self *Response, page string, parameters map[string]string) {
	if nil == parameters {
		parameters = map[string]string{}
	}

	p, pathFound := components[page]
	if !pathFound {
		NotifierSendError(self.server.notifier, fmt.Errorf("redirect to page `%s` failed because page id `%s` is unknown", page, page))
	}

	location := string(
		pathFieldRegex.ReplaceAllFunc(
			[]byte(p),
			func(i []byte) []byte {
				if nil == parameters {
					return []byte{}
				}
				key := string(i[1 : len(i)-1])
				return []byte(parameters[key])
			},
		),
	)

	self.navigate = &navigate{
		Page:       page,
		Parameters: parameters,
		Location:   location,
	}
}

// ResponseSendNavigate sends the client an instruction to navigate.
func ResponseSendNavigate(self *Response, page string) {
	ResponseSendNavigateWithParameters(self, page, map[string]string{})
}

// ResponseSendRedirect redirects the request.
func ResponseSendRedirect(self *Response, location string, statusCode int) {
	ResponseSendStatus(self, statusCode)
	ResponseSendHeader(self, "Location", location)
}

// ResponseSendRedirectToSecure redirects the request to the https server.
func ResponseSendRedirectToSecure(self *Response) {
	request := self.request
	if "" == request.server.certificate || "" == request.server.certificateKey || request.httpRequest.TLS != nil {
		return
	}

	insecureSuffix := fmt.Sprintf(":%d", request.server.port)
	secureSuffix := fmt.Sprintf(":%d", request.server.securePort)
	secureHost := strings.Replace(request.httpRequest.Host, insecureSuffix, secureSuffix, 1)
	secureLocation := fmt.Sprintf("https://%s%s", secureHost, request.httpRequest.RequestURI)
	ResponseSendRedirect(self, secureLocation, 302)
	return
}

// ResponseSendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func ResponseSendStatus(self *Response, code int) {
	if self.lockedStatusAndHeader {
		NotifierSendError(self.server.notifier, errors.New("status is locked"))
		return
	}
	self.statusCode = code
}

// ResponseSendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func ResponseSendHeader(self *Response, key string, value string) {
	if self.lockedStatusAndHeader {
		NotifierSendError(self.server.notifier, errors.New("headers locked"))
		return
	}

	self.header.Set(key, value)
}

// ResponseSendContentType sets the Content-Type header field.
func ResponseSendContentType(self *Response, contentType string) {
	ResponseSendHeader(self, "Content-Type", contentType)
}

// ResponseSendCookie sends a cookies to the client.
func ResponseSendCookie(self *Response, key string, value string) {
	ResponseSendHeader(self, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}

// ResponseSendContent sends binary safe content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func ResponseSendContent(self *Response, content []byte) {
	if !self.lockedStatusAndHeader {
		(*self.writer).WriteHeader(self.statusCode)
		self.lockedStatusAndHeader = true
	}

	if self.webSocket != nil {
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			NotifierSendError(self.server.notifier, writeError)
			return
		}
		return
	}

	if "" != self.eventName {
		sendEventContent(self, content)
		return
	}

	_, err := (*self.writer).Write(content)
	if err != nil {
		NotifierSendError(self.server.notifier, err)
		return
	}
}

// ResponseSendMessage sends utf-8 safe content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func ResponseSendMessage(self *Response, message string) {
	ResponseSendContent(self, []byte(message))
}

// ResponseSendNotFound sends an empty message with status 404 Not Found.
func ResponseSendNotFound(self *Response) {
	ResponseSendStatus(self, http.StatusNotFound)
}

// ResponseSendUnauthorized sends an empty message with status 401 Unauthorized.
func ResponseSendUnauthorized(self *Response) {
	ResponseSendStatus(self, http.StatusUnauthorized)
}

// ResponseSendBadRequest tris to send an empty message with status 400 Bad Request.
func ResponseSendBadRequest(self *Response) {
	ResponseSendStatus(self, http.StatusBadRequest)
}

// ResponseSendForbidden sends an empty message with status 403 Forbidden.
func ResponseSendForbidden(self *Response) {
	ResponseSendStatus(self, http.StatusForbidden)
}

// ResponseSendTooManyRequests sends and empty message with status 403 Forbidden.
func ResponseSendTooManyRequests(self *Response) {
	ResponseSendStatus(self, http.StatusTooManyRequests)
}

// ResponseSendJson sends json content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func ResponseSendJson(self *Response, payload any) {
	content, marshalError := json.Marshal(payload)
	if marshalError != nil {
		NotifierSendError(self.server.notifier, marshalError)
		return
	}

	if nil == self.webSocket {
		contentType := self.header.Get("Content-Type")
		if "" == contentType {
			self.header.Set("Content-Type", "application/json")
		}
	}

	ResponseSendContent(self, content)
}

// RequestVerifyContentType checks if the incoming request has any of the given content-types.
func RequestVerifyContentType(self *Request, contentTypes ...string) bool {
	requestedMime := self.httpRequest.Header.Get("Content-Type")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// RequestVerifyAccept checks if the incoming request accepts any of the given content-types.
func RequestVerifyAccept(self *Request, contentTypes ...string) bool {
	requestedAcceptMime := self.httpRequest.Header.Get("Accept")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.Contains(requestedAcceptMime, acceptedMime) {
			return true
		}
	}

	return false
}

// sendEventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func sendEventContent(self *Response, content []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", self.eventId, self.eventName)

	_, writeEventError := (*self.writer).Write([]byte(header))
	if writeEventError != nil {
		NotifierSendError(self.server.notifier, writeEventError)
		return
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeEventError = (*self.writer).Write([]byte("data: "))
		if writeEventError != nil {
			NotifierSendError(self.server.notifier, writeEventError)
			return
		}

		_, writeEventError = (*self.writer).Write(line)
		if writeEventError != nil {
			NotifierSendError(self.server.notifier, writeEventError)
			return
		}

		_, writeEventError = (*self.writer).Write([]byte("\r\n"))
		if writeEventError != nil {
			NotifierSendError(self.server.notifier, writeEventError)
			return
		}
	}

	_, writeEventError = (*self.writer).Write([]byte("\r\n"))
	if writeEventError != nil {
		NotifierSendError(self.server.notifier, writeEventError)
		return
	}

	flusher, flushedOk := (*self.writer).(http.Flusher)
	if !flushedOk {
		NotifierSendError(self.server.notifier, errors.New("could not retrieve flusher"))
		return
	}

	flusher.Flush()

	self.eventId++
}

// ResponseSendEmbeddedFileOrIndexOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func ResponseSendEmbeddedFileOrIndexOrElse(self *Response, orElse func()) {
	request := self.request
	fileName := filepath.Join(".dist", "client", request.httpRequest.RequestURI)

	if !existsInEmbeddedFileSystem(request.server.embeddedFileSystem, fileName) {
		orElse()
		return
	}

	if isEmbeddedDirectory(request.server.embeddedFileSystem, fileName) {
		fileName = filepath.Join(fileName, "index.html")
		if !isFile(fileName) {
			orElse()
			return
		}
	}

	reader, info, readerError := createReaderFromEmbeddedFileName(&request.server.embeddedFileSystem, fileName)
	if readerError != nil {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		sendEventContent(self, content)
		return
	}

	if "" == self.header.Get("Content-Type") {
		ResponseSendHeader(self, "Content-Type", mime(fileName))
	}

	if "" == self.header.Get("Content-Length") {
		ResponseSendHeader(self, "Content-Length", fmt.Sprintf("%d", (*info).Size()))
	}
	http.ServeContent(*self.writer, request.httpRequest, fileName, (*info).ModTime(), reader)
}

// ResponseSendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func ResponseSendEmbeddedFileOrElse(self *Response, orElse func()) {
	request := self.request
	fileName := filepath.Join(".dist", "client", request.httpRequest.RequestURI)
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !existsInEmbeddedFileSystem(request.server.embeddedFileSystem, fileName) ||
		isEmbeddedDirectory(request.server.embeddedFileSystem, fileName) {
		orElse()
		return
	}

	reader, info, readerError := createReaderFromEmbeddedFileName(&request.server.embeddedFileSystem, fileName)
	if readerError != nil {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		sendEventContent(self, content)
		return
	}

	if "" == self.header.Get("Content-Type") {
		ResponseSendHeader(self, "Content-Type", mime(fileName))
	}

	if "" == self.header.Get("Content-Length") {
		ResponseSendHeader(self, "Content-Length", fmt.Sprintf("%d", (*info).Size()))
	}
	http.ServeContent(*self.writer, request.httpRequest, fileName, (*info).ModTime(), reader)
}

// ResponseSendFileOrIndexOrElse sends the file requested by the client,
// or the closest index.html file, or else falls back.
func ResponseSendFileOrIndexOrElse(self *Response, orElse func()) {
	request := self.request
	fileName := filepath.Join(".dist", "client", request.httpRequest.RequestURI)

	if !exists(fileName) {
		orElse()
		return
	}

	if isDirectory(fileName) {
		fileName = filepath.Join(fileName, "index.html")
		if !isFile(fileName) {
			orElse()
			return
		}
	}

	reader, info, readerError := createReaderFromFileName(fileName)
	if readerError != nil {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		sendEventContent(self, content)
		return
	}

	if "" == self.header.Get("Content-Type") {
		ResponseSendHeader(self, "Content-Type", mime(fileName))
	}

	if "" == self.header.Get("Content-Length") {
		ResponseSendHeader(self, "Content-Length", fmt.Sprintf("%d", (*info).Size()))
	}
	http.ServeContent(*self.writer, request.httpRequest, fileName, (*info).ModTime(), reader)
}

// ResponseSendFileOrElse sends the file requested by the client, or else falls back.
func ResponseSendFileOrElse(self *Response, orElse func()) {
	request := self.request
	fileName := filepath.Join(".dist", "client", request.httpRequest.RequestURI)

	if !exists(fileName) || isDirectory(fileName) {
		orElse()
		return
	}

	reader, info, readerError := createReaderFromFileName(fileName)
	if readerError != nil {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		sendEventContent(self, content)
		return
	}

	if "" == self.header.Get("Content-Type") {
		ResponseSendHeader(self, "Content-Type", mime(fileName))
	}

	if "" == self.header.Get("Content-Length") {
		ResponseSendHeader(self, "Content-Length", fmt.Sprintf("%d", (*info).Size()))
	}
	http.ServeContent(*self.writer, request.httpRequest, fileName, (*info).ModTime(), reader)
}

func createReaderFromEmbeddedFileName(efs *embed.FS, fileName string) (*bytes.Reader, *os.FileInfo, error) {
	file, openError := efs.Open(fileName)
	if openError != nil {
		return nil, nil, openError
	}

	fileInfo, _ := file.Stat()

	buffer := make([]byte, fileInfo.Size())
	_, readError := file.Read(buffer)
	if readError != nil {
		closeError := file.Close()
		if closeError != nil {
			return nil, nil, closeError
		}
		return nil, nil, readError
	}

	closeError := file.Close()
	if closeError != nil {
		return nil, nil, closeError
	}
	return bytes.NewReader(buffer), &fileInfo, nil
}

func createReaderFromFileName(fileName string) (*bytes.Reader, *os.FileInfo, error) {
	file, openError := os.Open(fileName)
	if openError != nil {
		return nil, nil, openError
	}

	fileInfo, _ := file.Stat()

	buffer := make([]byte, fileInfo.Size())
	_, readError := file.Read(buffer)
	if readError != nil {
		closeError := file.Close()
		if closeError != nil {
			return nil, nil, closeError
		}
		return nil, nil, readError
	}

	closeError := file.Close()
	if closeError != nil {
		return nil, nil, closeError
	}
	return bytes.NewReader(buffer), &fileInfo, nil
}

// ResponseSendSseUpgrade upgrades the http connection to server sent events
// and returns a function that sets the name of the current event.
//
// The default event is "message".
func ResponseSendSseUpgrade(self *Response) (setEventName func(eventName string)) {
	ResponseSendHeader(self, "Access-Control-Allow-Origin", "*")
	ResponseSendHeader(self, "Access-Control-Expose-Headers", "Content-Type")
	ResponseSendHeader(self, "Content-Type", "text/event-stream")
	ResponseSendHeader(self, "Cache-Control", "no-cache")
	ResponseSendHeader(self, "Connection", "keep-alive")
	self.eventName = "message"
	setEventName = func(eventName string) {
		if "" == eventName {
			NotifierSendError(
				self.server.notifier,
				fmt.Errorf("renaming a server sent event (`%s`) to an empty string is not allowed", self.eventName),
			)
			return
		}

		self.eventName = eventName
	}
	return
}

// ResponseSendWsUpgrade upgrades the http connection to web sockets.
func ResponseSendWsUpgrade(self *Response) {
	request := self.request
	conn, upgradeError := self.server.webSocketUpgrader.Upgrade(*self.writer, request.httpRequest, nil)
	if upgradeError != nil {
		NotifierSendError(request.server.notifier, upgradeError)
		return
	}
	defer func(conn *websocket.Conn) {
		closeError := conn.Close()
		if closeError != nil {
			NotifierSendError(request.server.notifier, closeError)
		}
	}(conn)
	self.webSocket = conn
	request.webSocketConn = conn
	self.lockedStatusAndHeader = true
}

// ResponseSendView sends a view.
func ResponseSendView(self *Response, view *View) {
	embeddedFileSystem := view.embeddedFileSystem
	if nil == embeddedFileSystem {
		embeddedFileSystem = &self.server.embeddedFileSystem
	}

	notifier := view.notifier
	if nil == embeddedFileSystem {
		notifier = self.server.notifier
	}

	content, compileError := ViewRender(&View{
		render:             view.render,
		data:               view.data,
		name:               view.name,
		parameters:         view.parameters,
		functions:          view.functions,
		embeddedFileSystem: embeddedFileSystem,
		notifier:           notifier,
	})
	if nil != compileError {
		NotifierSendError(self.server.notifier, compileError)
		return
	}

	if "" == self.header.Get("Content-Type") {
		ResponseSendHeader(self, "Content-Type", "text/html")
	}

	ResponseSendMessage(self, content)
}

type Api struct {
	patterns []string
	handler  func(request *Request, response *Response)
	guards   []func(request *Request, response *Response, pass func())
}

// ApiBuilder builds an api.
type ApiBuilder = func(api *Api)

// ApiWithPattern adds a pattern.
func ApiWithPattern(self *Api, pattern string) {
	self.patterns = append(self.patterns, pattern)
}

// ApiWithRequestHandler sets the request handler.
func ApiWithRequestHandler(self *Api, handler func(request *Request, response *Response)) {
	self.handler = handler
}

// ApiWithGuardHandler add a guard handler.
func ApiWithGuardHandler(self *Api, handler func(request *Request, response *Response, pass func())) {
	self.guards = append(self.guards, handler)
}

// ServerWithApiBuilder adds an api.
func ServerWithApiBuilder(self *Server, builder ApiBuilder) {
	api := &Api{
		patterns: []string{},
		guards:   []func(request *Request, response *Response, pass func()){},
	}

	builder(api)

	if nil == api.handler {
		api.handler = func(request *Request, response *Response) {
			// Noop.
		}
	}

	for _, pattern := range api.patterns {
		if "" == pattern {
			NotifierSendError(self.notifier, fmt.Errorf("could not add api because path is empty"))
			return
		}
		serverMapRoute(self, pattern, routeCreate(api.handler, api.guards))
	}
}

type Page struct {
	paths  []string
	view   *View
	base   func(request *Request, response *Response, view *View)
	action func(request *Request, response *Response, view *View)
	guards []func(request *Request, response *Response, pass func())
}

// PageBuilder builds a page.
type PageBuilder = func(page *Page)

// PageWithPath adds a path.
func PageWithPath(self *Page, path string) {
	self.paths = append(self.paths, path)
}

// PageWithView sets the view.
func PageWithView(self *Page, view *View) {
	self.view = view
}

// PageWithBaseHandler sets the base.
func PageWithBaseHandler(self *Page, handler func(request *Request, response *Response, view *View)) {
	self.base = handler
}

// PageWithActionHandler sets the action handler.
func PageWithActionHandler(self *Page, handler func(request *Request, response *Response, view *View)) {
	self.action = handler
}

// PageWithGuardHandler add a guard handler.
func PageWithGuardHandler(self *Page, handler func(request *Request, response *Response, pass func())) {
	self.guards = append(self.guards, handler)
}

// ServerWithPageBuilder adds a page.
func ServerWithPageBuilder(self *Server, builder PageBuilder) {
	page := &Page{
		view:   ViewReference("Default"),
		paths:  []string{},
		guards: []func(request *Request, response *Response, pass func()){},
	}

	builder(page)

	if 0 == len(page.paths) {
		page.paths = append(page.paths, "/"+strings.ReplaceAll(page.view.name, ".", "/"))
	}

	if "" == page.view.name {
		NotifierSendError(self.notifier, fmt.Errorf("view name cannot be empty"))
		return
	}

	if nil == page.base {
		page.base = func(request *Request, response *Response, view *View) {
			// Noop.
		}
	}

	if nil == page.action {
		page.action = func(request *Request, response *Response, view *View) {
			// Noop.
		}
	}

	for _, path_ := range page.paths {
		serverMapRoute(self, "GET "+path_, routeCreateWithView(page.base, page.guards, page.view.name))
		serverMapRoute(self, "POST "+path_, routeCreateWithView(page.action, page.guards, page.view.name))
	}
}
