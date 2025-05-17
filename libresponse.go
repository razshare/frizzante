package frizzante

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

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
	context               map[string]any
}

type navigate struct {
	Page       string
	Parameters map[string]string
	Location   string
}

var pathFieldRegex = regexp.MustCompile(`\{(.*?)}`)

// responseSendEventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func responseSendEventContent(self *Response, content []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", self.eventId, self.eventName)

	_, writeEventError := (*self.writer).Write([]byte(header))
	if nil != writeEventError {
		NotifierSendError(self.server.notifier, writeEventError)
		return
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeEventError = (*self.writer).Write([]byte("data: "))
		if nil != writeEventError {
			NotifierSendError(self.server.notifier, writeEventError)
			return
		}

		_, writeEventError = (*self.writer).Write(line)
		if nil != writeEventError {
			NotifierSendError(self.server.notifier, writeEventError)
			return
		}

		_, writeEventError = (*self.writer).Write([]byte("\r\n"))
		if nil != writeEventError {
			NotifierSendError(self.server.notifier, writeEventError)
			return
		}
	}

	_, writeEventError = (*self.writer).Write([]byte("\r\n"))
	if nil != writeEventError {
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

// ResponseSendNavigateWithParameters sends the client an instruction to navigate.
func ResponseSendNavigateWithParameters(self *Response, page string, parameters map[string]string) {
	if nil == parameters {
		parameters = map[string]string{}
	}

	self.navigate = &navigate{
		Page:       page,
		Parameters: parameters,
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
		if nil != writeError {
			NotifierSendError(self.server.notifier, writeError)
			return
		}
		return
	}

	if "" != self.eventName {
		responseSendEventContent(self, content)
		return
	}

	_, err := (*self.writer).Write(content)
	if nil != err {
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

// ResponseSendBadRequest sends an empty message with status 400 Bad Request.
func ResponseSendBadRequest(self *Response) {
	ResponseSendStatus(self, http.StatusBadRequest)
}

// ResponseSendInternalServerError sends an error message with status 500 Internal Server Error
// and also sends the error to the server notifier.
func ResponseSendInternalServerError(self *Response, err error) {
	NotifierSendError(self.server.notifier, err)
	ResponseSendStatus(self, http.StatusBadRequest)
	ResponseSendMessage(self, err.Error())
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
	if nil != marshalError {
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

// ResponseSendEmbeddedFileOrIndexOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func ResponseSendEmbeddedFileOrIndexOrElse(self *Response, orElse func()) {
	hasEmbeddedFileSystem := nil != self.request.server.embeddedFileSystem
	if !hasEmbeddedFileSystem {
		orElse()
		return
	}

	request := self.request
	fileName := filepath.Join(".dist", "client", request.httpRequest.RequestURI)

	if !existsInEmbeddedFileSystem(*request.server.embeddedFileSystem, fileName) {
		orElse()
		return
	}

	if isEmbeddedDirectory(*request.server.embeddedFileSystem, fileName) {
		fileName = filepath.Join(fileName, "index.html")
		if !isFile(fileName) {
			orElse()
			return
		}
	}

	reader, info, readerError := createReaderFromEmbeddedFileName(request.server.embeddedFileSystem, fileName)
	if nil != readerError {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		responseSendEventContent(self, content)
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
	hasEmbeddedFileSystem := nil != self.request.server.embeddedFileSystem
	if !hasEmbeddedFileSystem {
		orElse()
		return
	}

	request := self.request
	fileName := filepath.Join(".dist", "client", request.httpRequest.RequestURI)
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !existsInEmbeddedFileSystem(*request.server.embeddedFileSystem, fileName) ||
		isEmbeddedDirectory(*request.server.embeddedFileSystem, fileName) {
		orElse()
		return
	}

	reader, info, readerError := createReaderFromEmbeddedFileName(request.server.embeddedFileSystem, fileName)
	if nil != readerError {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		responseSendEventContent(self, content)
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

	if !fileExists(fileName) {
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
	if nil != readerError {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		responseSendEventContent(self, content)
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

	if !fileExists(fileName) || isDirectory(fileName) {
		orElse()
		return
	}

	reader, info, readerError := createReaderFromFileName(fileName)
	if nil != readerError {
		NotifierSendError(self.server.notifier, readerError)
		return
	}

	if self.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		writeError := self.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			NotifierSendError(self.server.notifier, writeError)
		}
		return
	}

	if "" != self.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return
		}
		responseSendEventContent(self, content)
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
	if nil != upgradeError {
		NotifierSendError(request.server.notifier, upgradeError)
		return
	}
	defer func(conn *websocket.Conn) {
		closeError := conn.Close()
		if nil != closeError {
			NotifierSendError(request.server.notifier, closeError)
		}
	}(conn)
	self.webSocket = conn
	request.webSocketConn = conn
	self.lockedStatusAndHeader = true
}

// ResponseSendView sends a view.
func ResponseSendView[T any](self *Response, view *View[T]) {
	content, compileError := ViewRender(view)
	if nil != compileError {
		NotifierSendError(self.server.notifier, compileError)
		return
	}

	if "" == self.header.Get("Content-Type") {
		ResponseSendHeader(self, "Content-Type", "text/html")
	}

	ResponseSendMessage(self, content)
}
