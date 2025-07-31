package actions

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/mimes"
	"github.com/razshare/frizzante/traces"
	"github.com/razshare/frizzante/views"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
////////////////////////// RECEIVE //////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////

// ReceiveCancellation returns a channel that closes when the request gets cancelled.
func ReceiveCancellation(self *connections.Connection) <-chan struct{} {
	return self.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(self *connections.Connection) *bool {
	isAlive := true
	go func() {
		<-ReceiveCancellation(self)
		isAlive = false
	}()
	return &isAlive
}

// ReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func ReceiveCookie(self *connections.Connection, key string) string {
	cookie, cookieError := self.Request.Cookie(key)
	if cookieError != nil {
		traces.Trace(self.Http.ErrorLog, cookieError)
		return ""
	}

	data, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		traces.Trace(self.Http.ErrorLog, queryError)
		return ""
	}

	return data
}

// ReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func ReceiveMessage(self *connections.Connection) string {
	if self.WebSocket != nil {
		_, data, readError := self.WebSocket.ReadMessage()
		if readError != nil {
			traces.Trace(self.Http.ErrorLog, readError)
			return ""
		}
		return string(data)
	}

	data, readError := io.ReadAll(self.Request.Body)
	if readError != nil {
		traces.Trace(self.Http.ErrorLog, readError)
		return ""
	}
	return string(data)
}

// ReceiveJson reads the next JSON-encoded message from the
// connection and stores it in the value pointed to by val.
//
// Compatible with web sockets.
func ReceiveJson(self *connections.Connection, value any) {
	if self.WebSocket != nil {
		jsonError := self.WebSocket.ReadJSON(value)
		if jsonError != nil {
			traces.Trace(self.Http.ErrorLog, jsonError)
			return
		}
		return
	}

	data, readError := io.ReadAll(self.Request.Body)
	if readError != nil {
		traces.Trace(self.Http.ErrorLog, readError)
		return
	}

	jsonError := json.Unmarshal(data, value)
	if jsonError != nil {
		traces.Trace(self.Http.ErrorLog, jsonError)
		return
	}
}

// ReceiveForm reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of 2MB
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func ReceiveForm(self *connections.Connection) url.Values {
	return ReceiveFormWithMaxMemory(self, 2*globals.MB)
}

// ReceiveFormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of maxMemory bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func ReceiveFormWithMaxMemory(self *connections.Connection, maxMemory int64) url.Values {
	if self.WebSocket != nil {
		traces.Trace(self.Http.ErrorLog, errors.New("connection is not of type web socket"))
		return url.Values{}
	}

	formError := self.Request.ParseMultipartForm(maxMemory)
	if formError != nil {
		if !errors.Is(formError, http.ErrNotMultipart) {
			return url.Values{}
		}

		formError = self.Request.ParseForm()
		if formError != nil {
			traces.Trace(self.Http.ErrorLog, formError)
			return url.Values{}
		}
	}

	return self.Request.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func ReceiveQuery(self *connections.Connection, key string) string {
	return self.Request.URL.Query().Get(key)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func ReceivePath(self *connections.Connection, key string) string {
	return self.Request.PathValue(key)
}

// ReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func ReceiveHeader(self *connections.Connection, key string) string {
	return self.Request.Header.Get(key)
}

// ReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ReceiveContentType(self *connections.Connection) string {
	return self.Request.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func VerifyContentType(self *connections.Connection, contentType ...string) bool {
	requestedMime := self.Request.Header.Get("Content-Type")
	for _, acceptedMime := range contentType {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func VerifyAccept(self *connections.Connection, accept ...string) bool {
	requestedAcceptMime := self.Request.Header.Get("Accept")
	for _, acceptedMime := range accept {
		if acceptedMime == "*" || strings.Contains(requestedAcceptMime, acceptedMime) {
			return true
		}
	}

	return false
}

/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
//////////////////////////// SEND ///////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////

// SendEventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a Server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func SendEventContent(self *connections.Connection, content []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", self.EventId, self.EventName)

	_, writeError := self.Writer.Write([]byte(header))
	if writeError != nil {
		traces.Trace(self.Http.ErrorLog, writeError)
		return
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeError = self.Writer.Write([]byte("data: "))
		if writeError != nil {
			traces.Trace(self.Http.ErrorLog, writeError)
			return
		}

		_, writeError = self.Writer.Write(line)
		if writeError != nil {
			traces.Trace(self.Http.ErrorLog, writeError)
			return
		}

		_, writeError = self.Writer.Write([]byte("\r\n"))
		if writeError != nil {
			traces.Trace(self.Http.ErrorLog, writeError)
			return
		}
	}

	_, writeError = self.Writer.Write([]byte("\r\n"))
	if writeError != nil {
		traces.Trace(self.Http.ErrorLog, writeError)
		return
	}

	flusher, flushedOk := self.Writer.(http.Flusher)
	if !flushedOk {
		traces.Trace(self.Http.ErrorLog, errors.New("could not retrieve flusher"))
		return
	}

	flusher.Flush()

	self.EventId++
}

// SendNavigate redirects the request to a location with status 302.
func SendNavigate(self *connections.Connection, location string) {
	SendRedirect(self, location, 302)
	SendFlush(self)
}

// SendRedirect redirects the request to a location with a status.
func SendRedirect(self *connections.Connection, location string, status int) {
	SendStatus(self, status)
	SendHeader(self, "Location", location)
}

// SendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func SendStatus(self *connections.Connection, status int) {
	if self.Locked {
		traces.Trace(self.Http.ErrorLog, "status is locked")
		return
	}

	self.Status = status
}

// SendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func SendHeader(self *connections.Connection, key string, value string) {
	if self.Locked {
		traces.Trace(self.Http.ErrorLog, "header is locked")
		return
	}

	self.Writer.Header().Set(key, value)
}

func SendHeaders(self *connections.Connection, headers map[string]string) {
	if self.Locked {
		traces.Trace(self.Http.ErrorLog, "header is locked")
		return
	}

	for key, value := range headers {
		self.Writer.Header().Set(key, value)
	}
}

// SendContentType sets the Content-Type header field.
func SendContentType(self *connections.Connection, contentType string) {
	SendHeader(self, "Content-Type", contentType)
}

// SendCookie sends a cookies to the client.
func SendCookie(self *connections.Connection, key string, value string) {
	SendHeader(self, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}

func SendFlush(self *connections.Connection) {
	SendMessage(self, "")
}

// SendContent sends binary safe content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func SendContent(self *connections.Connection, content []byte) {
	if !self.Locked {
		self.Writer.WriteHeader(self.Status)
		self.Locked = true
	}

	if self.WebSocket != nil {
		writeError := self.WebSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			traces.Trace(self.Http.ErrorLog, writeError)
		}
		return
	}

	if "" != self.EventName {
		SendEventContent(self, content)
		return
	}

	_, writeError := self.Writer.Write(content)
	if writeError != nil {
		traces.Trace(self.Http.ErrorLog, writeError)
	}
}

// SendMessage sends utf-8 safe content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func SendMessage(self *connections.Connection, value string) {
	SendContent(self, []byte(value))
}

// SendNotFound sends a message with status 404 Not Found.
func SendNotFound(self *connections.Connection, value string) {
	SendStatus(self, http.StatusNotFound)
	SendMessage(self, value)
}

// SendUnauthorized sends a message with status 401 Unauthorized.
func SendUnauthorized(self *connections.Connection, value string) {
	SendStatus(self, http.StatusUnauthorized)
	SendMessage(self, value)
}

// SendBadRequest sends a message with status 400 Bad Request.
func SendBadRequest(self *connections.Connection, value string) {
	SendStatus(self, http.StatusBadRequest)
	SendMessage(self, value)
}

// SendError sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func SendError(self *connections.Connection, err error) {
	SendStatus(self, http.StatusBadRequest)
	SendMessage(self, err.Error())
}

// SendForbidden sends a message with status 403 Forbidden.
func SendForbidden(self *connections.Connection, v string) {
	SendStatus(self, http.StatusForbidden)
	SendMessage(self, v)
}

// SendTooManyRequests sends a message with status 403 Forbidden.
func SendTooManyRequests(self *connections.Connection, v string) {
	SendStatus(self, http.StatusTooManyRequests)
	SendMessage(self, v)
}

// SendJson sends json content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func SendJson(self *connections.Connection, value any) {
	data, jsonError := json.Marshal(value)
	if jsonError != nil {
		traces.Trace(self.Http.ErrorLog, jsonError)
		return
	}

	if nil == self.WebSocket {
		contentType := self.Writer.Header().Get("Content-Type")
		if "" == contentType {
			self.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	SendContent(self, data)
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func SendEmbeddedFileOrElse(self *connections.Connection, efs embed.FS, orElse func()) {
	fileName := self.PublicRoot + self.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !embeds.IsFile(efs, fileName) || embeds.IsDirectory(efs, fileName) {
		orElse()
		return
	}

	reader, readerInfo, readerError := embeds.FileReader(efs, fileName)
	if readerError != nil {
		traces.Trace(self.Http.ErrorLog, readerError)
		return
	}

	if self.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(self.Http.ErrorLog, readError)
			return
		}

		writeError := self.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			traces.Trace(self.Http.ErrorLog, writeError)
			return
		}
		return
	}

	if "" != self.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(self.Http.ErrorLog, readError)
			return
		}

		SendEventContent(self, data)
		return
	}

	if "" == self.Writer.Header().Get("Content-Type") {
		SendHeader(self, "Content-Type", mimes.Mime(fileName))
	}

	if "" == self.Writer.Header().Get("Content-Length") {
		SendHeader(self, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(self.Writer, self.Request, fileName, readerInfo.ModTime(), reader)
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func SendFileOrElse(self *connections.Connection, orElse func()) {
	fileName := filepath.Join(self.PublicRoot, self.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		SendEmbeddedFileOrElse(self, self.Efs, orElse)
		return
	}

	reader, readerInfo, readerError := files.FileReader(fileName)
	if readerError != nil {
		traces.Trace(self.Http.ErrorLog, readerError)
		return
	}

	if self.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(self.Http.ErrorLog, readError)
			return
		}

		writeError := self.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			traces.Trace(self.Http.ErrorLog, writeError)
			return
		}
	}

	if "" != self.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(self.Http.ErrorLog, readError)
			return
		}

		SendEventContent(self, data)
		return
	}

	if "" == self.Writer.Header().Get("Content-Type") {
		SendHeader(self, "Content-Type", mimes.Mime(fileName))
	}

	if "" == self.Writer.Header().Get("Content-Length") {
		SendHeader(self, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(self.Writer, self.Request, fileName, readerInfo.ModTime(), reader)
}

// SendSseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SendSseUpgrade(self *connections.Connection) func(eventName string) {
	SendHeaders(self, map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Connection":                    "keep-alive",
	})

	self.EventName = "message"

	return func(eventName string) { self.EventName = eventName }
}

// SendWsUpgrade upgrades to web sockets.
func SendWsUpgrade(self *connections.Connection) {
	SendConfiguredWsUpgrade(self, websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// SendConfiguredWsUpgrade upgrades to web sockets.
func SendConfiguredWsUpgrade(self *connections.Connection, upgrader websocket.Upgrader) {
	webSocketConnection, upgradeError := upgrader.Upgrade(self.Writer, self.Request, nil)
	if upgradeError != nil {
		traces.Trace(self.Http.ErrorLog, upgradeError)
		return
	}

	defer func(webSocketConnection *websocket.Conn) {
		closeError := webSocketConnection.Close()
		if closeError != nil {
			traces.Trace(self.Http.ErrorLog, closeError)
		}
	}(webSocketConnection)

	self.WebSocket = webSocketConnection
	self.Locked = true

	return
}

// SendView sends a view.
func SendView(self *connections.Connection, view views.View) {
	if self.Writer.Header().Get("Location") != "" {
		return
	}

	if VerifyAccept(self, "application/json") {
		if view.Data == nil {
			view.Data = map[string]any{}
		}
		props := map[string]any{
			"name":       view.Name,
			"data":       view.Data,
			"renderMode": view.RenderMode,
		}
		SendJson(self, props)
		return
	}

	if view.ServerJs == "" {
		view.ServerJs = self.ServerJs
	}

	if view.IndexHtml == "" {
		view.IndexHtml = self.IndexHtml
	}

	if view.AppRoot == "" {
		view.AppRoot = self.AppRoot
	}

	html, renderError := views.Render(&view, self.Efs)
	if renderError != nil {
		traces.Trace(self.Http.ErrorLog, renderError)
	}

	if "" == self.Writer.Header().Get("Content-Type") {
		SendHeader(self, "Content-Type", "text/html")
	}

	SendMessage(self, html)
}
