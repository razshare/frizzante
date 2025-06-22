package frz

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type Connection struct {
	server    *Server
	request   *http.Request
	writer    http.ResponseWriter
	locked    bool
	status    int
	header    http.Header
	webSocket *websocket.Conn
	eventName string
	eventId   int64
	sessionId string
}

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
func (connection *Connection) ReceiveCancellation() <-chan struct{} {
	return connection.request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func (connection *Connection) IsAlive() *bool {
	value := true
	go func() {
		<-connection.ReceiveCancellation()
		value = false
	}()
	return &value
}

// ReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// It silently discards malformed values.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveCookie(key string) string {
	cookie, cookieError := connection.request.Cookie(key)
	if nil != cookieError {
		return ""
	}
	value, unescapeError := url.QueryUnescape(cookie.Value)
	if nil != unescapeError {
		return ""
	}

	return value
}

// ReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveMessage() string {
	if connection.webSocket != nil {
		_, readBytes, readError := connection.webSocket.ReadMessage()
		if nil != readError {
			connection.server.notifier.SendErrorAndTrace(readError, 1)
			return ""
		}
		return string(readBytes)
	}

	readBytes, readAllError := io.ReadAll(connection.request.Body)
	if nil != readAllError {
		connection.server.notifier.SendErrorAndTrace(readAllError, 1)
		return ""
	}
	return string(readBytes)
}

// ReceiveJson reads the next JSON-encoded message from the
// connection and stores it in the value pointed to by v.
//
// ReceiveJson returns true on success or false on failure.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveJson(v any) bool {
	if connection.webSocket != nil {
		jsonError := connection.webSocket.ReadJSON(v)
		if nil != jsonError {
			connection.server.notifier.SendErrorAndTrace(jsonError, 1)
			return false
		}
		return true
	}

	readBytes, readAllError := io.ReadAll(connection.request.Body)
	if nil != readAllError {
		connection.server.notifier.SendErrorAndTrace(readAllError, 1)
		return false
	}
	unmarshalError := json.Unmarshal(readBytes, v)
	if nil != unmarshalError {
		connection.server.notifier.SendErrorAndTrace(unmarshalError, 1)
		return false
	}
	return true
}

// ReceiveForm reads the message as a form and returns the value.
//
// It silently discards malformed values.
func (connection *Connection) ReceiveForm() url.Values {
	return connection.ReceiveFormWithMaxMemory(2 * MB)
}

// ReceiveFormWithMaxMemory reads the message as a form and returns the value.
//
// It silently discards malformed values.
func (connection *Connection) ReceiveFormWithMaxMemory(maxMemory int64) url.Values {
	if connection.webSocket != nil {
		return url.Values{}
	}

	parseMultipartFormError := connection.request.ParseMultipartForm(maxMemory)
	if nil != parseMultipartFormError {
		if !errors.Is(parseMultipartFormError, http.ErrNotMultipart) {
			return url.Values{}
		}

		parseFormError := connection.request.ParseForm()
		if nil != parseFormError {
			return url.Values{}
		}
	}

	return connection.request.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveQuery(name string) string {
	return connection.request.URL.Query().Get(name)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceivePath(name string) string {
	return connection.request.PathValue(name)
}

// ReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveHeader(key string) string {
	return connection.request.Header.Get(key)
}

// ReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ReceiveContentType(connection *Connection) string {
	return connection.request.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func VerifyContentType(connection *Connection, contentTypes ...string) bool {
	requestedMime := connection.request.Header.Get("Content-Type")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func (connection *Connection) VerifyAccept(contentTypes ...string) bool {
	requestedAcceptMime := connection.request.Header.Get("Accept")
	for _, acceptedMime := range contentTypes {
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
func (connection *Connection) SendEventContent(content []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", connection.eventId, connection.eventName)

	_, writeEventError := connection.writer.Write([]byte(header))
	if nil != writeEventError {
		connection.server.notifier.SendErrorAndTrace(writeEventError, 1)
		return
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeEventError = connection.writer.Write([]byte("data: "))
		if nil != writeEventError {
			connection.server.notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}

		_, writeEventError = connection.writer.Write(line)
		if nil != writeEventError {
			connection.server.notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}

		_, writeEventError = connection.writer.Write([]byte("\r\n"))
		if nil != writeEventError {
			connection.server.notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}
	}

	_, writeEventError = connection.writer.Write([]byte("\r\n"))
	if nil != writeEventError {
		connection.server.notifier.SendErrorAndTrace(writeEventError, 1)
		return
	}

	flusher, flushedOk := connection.writer.(http.Flusher)
	if !flushedOk {
		connection.server.notifier.SendErrorAndTrace(errors.New("could not retrieve flusher"), 1)
		return
	}

	flusher.Flush()

	connection.eventId++
}

// SendNavigate redirects the request with status 302.
func (connection *Connection) SendNavigate(location string) {
	connection.SendRedirect(location, 302)
	connection.SendFlush()
}

// SendRedirect redirects the request.
func (connection *Connection) SendRedirect(location string, statusCode int) {
	connection.SendStatus(statusCode)
	connection.SendHeader("Location", location)
}

// SendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func (connection *Connection) SendStatus(code int) {
	if connection.locked {
		connection.server.notifier.SendErrorAndTrace(errors.New("status is locked"), 1)
	}
	connection.status = code
}

// SendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func (connection *Connection) SendHeader(key string, value string) {
	if connection.locked {
		connection.server.notifier.SendErrorAndTrace(errors.New("header is locked"), 1)
	}

	connection.header.Set(key, value)
}

// SendContentType sets the Content-Type header field.
func (connection *Connection) SendContentType(contentType string) {
	connection.SendHeader("Content-Type", contentType)
}

// SendCookie sends a cookies to the client.
func (connection *Connection) SendCookie(key string, value string) {
	connection.SendHeader(
		"Set-Cookie",
		fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)),
	)
}

func (connection *Connection) SendFlush() {
	connection.SendMessage("")
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
func (connection *Connection) SendContent(content []byte) {
	if !connection.locked {
		connection.writer.WriteHeader(connection.status)
		connection.locked = true
	}

	if connection.webSocket != nil {
		writeError := connection.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			connection.server.notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != connection.eventName {
		connection.SendEventContent(content)
		return
	}

	_, writeError := connection.writer.Write(content)
	if nil != writeError {
		connection.server.notifier.SendErrorAndTrace(writeError, 1)
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
func (connection *Connection) SendMessage(message string) {
	connection.SendContent([]byte(message))
}

// SendNotFound sends a message with status 404 Not Found.
func (connection *Connection) SendNotFound(message string) {
	connection.SendStatus(http.StatusNotFound)
	connection.SendMessage(message)
}

// SendUnauthorized sends a message with status 401 Unauthorized.
func (connection *Connection) SendUnauthorized(message string) {
	connection.SendStatus(http.StatusUnauthorized)
	connection.SendMessage(message)
}

// SendBadRequest sends a message with status 400 Bad Request.
func (connection *Connection) SendBadRequest(message string) {
	connection.SendStatus(http.StatusBadRequest)
	connection.SendMessage(message)
}

// SendInternalServerError sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func (connection *Connection) SendInternalServerError(err error) {
	connection.SendStatus(http.StatusBadRequest)
	connection.SendMessage(err.Error())
}

// SendForbidden sends a message with status 403 Forbidden.
func (connection *Connection) SendForbidden(message string) {
	connection.SendStatus(http.StatusForbidden)
	connection.SendMessage(message)
}

// SendTooManyRequests sends a message with status 403 Forbidden.
func (connection *Connection) SendTooManyRequests(message string) {
	connection.SendStatus(http.StatusTooManyRequests)
	connection.SendMessage(message)
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
func (connection *Connection) SendJson(payload any) {
	content, marshalError := json.Marshal(payload)
	if nil != marshalError {
		connection.server.notifier.SendErrorAndTrace(marshalError, 1)
		return
	}

	if nil == connection.webSocket {
		contentType := connection.header.Get("Content-Type")
		if "" == contentType {
			connection.header.Set("Content-Type", "application/json")
		}
	}

	connection.SendContent(content)
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func (connection *Connection) SendEmbeddedFileOrElse(efs embed.FS, orElse func()) {
	fileName := connection.server.publicRoot + connection.request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !EfsIsFile(efs, fileName) || EfsIsDirectory(efs, fileName) {
		orElse()
		return
	}

	reader, info, readerError := EfsFileReader(efs, fileName)
	if nil != readerError {
		connection.server.notifier.SendErrorAndTrace(readerError, 1)
		return
	}

	if connection.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.server.notifier.SendErrorAndTrace(readError, 1)
			return
		}
		writeError := connection.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			connection.server.notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != connection.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.server.notifier.SendErrorAndTrace(readError, 1)
		}
		connection.SendEventContent(content)
		return
	}

	if "" == connection.header.Get("Content-Type") {
		connection.SendHeader("Content-Type", Mime(fileName))
	}

	if "" == connection.header.Get("Content-Length") {
		connection.SendHeader("Content-Length", fmt.Sprintf("%d", info.Size()))
	}
	http.ServeContent(connection.writer, connection.request, fileName, info.ModTime(), reader)
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func (connection *Connection) SendFileOrElse(orElse func()) {
	fileName := filepath.Join(connection.server.publicRoot, connection.request.RequestURI)

	if !IsFile(fileName) || IsDirectory(fileName) {
		connection.SendEmbeddedFileOrElse(connection.server.efs, orElse)
		return
	}

	reader, info, readerError := FileReader(fileName)
	if nil != readerError {
		connection.server.notifier.SendErrorAndTrace(readerError, 1)
		return
	}

	if connection.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.server.notifier.SendErrorAndTrace(readError, 1)
			return
		}
		writeError := connection.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			connection.server.notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != connection.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.server.notifier.SendErrorAndTrace(readError, 1)
			return
		}
		connection.SendEventContent(content)
	}

	if "" == connection.header.Get("Content-Type") {
		connection.SendHeader("Content-Type", Mime(fileName))
	}

	if "" == connection.header.Get("Content-Length") {
		connection.SendHeader("Content-Length", fmt.Sprintf("%d", info.Size()))
	}
	http.ServeContent(connection.writer, connection.request, fileName, info.ModTime(), reader)
}

// SendSseUpgrade upgrades the http connection to server sent events
// and returns a function that sets the name of the current event.
//
// The default event is "message".
func SendSseUpgrade(connection *Connection) func(eventName string) {
	connection.SendHeader("Access-Control-Allow-Origin", "*")
	connection.SendHeader("Access-Control-Expose-Headers", "Content-Type")
	connection.SendHeader("Content-Type", "text/event-stream")
	connection.SendHeader("MemoryCache-Control", "no-cache")
	connection.SendHeader("Connection", "keep-alive")
	connection.eventName = "message"
	return func(eventName string) {
		if "" == eventName {
			connection.server.notifier.SendErrorAndTrace(
				fmt.Errorf("renaming a server sent event (`%s`) to an empty string is not allowed", connection.eventName),
				1,
			)
		}
		connection.eventName = eventName
	}
}

// SendSseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func (connection *Connection) SendSseUpgrade() (setEventName func(eventName string)) {
	connection.SendHeader("Access-Control-Allow-Origin", "*")
	connection.SendHeader("Access-Control-Expose-Headers", "Content-Type")
	connection.SendHeader("Content-Type", "text/event-stream")
	connection.SendHeader("Cache-Control", "no-cache")
	connection.SendHeader("Connection", "keep-alive")
	connection.eventName = "message"
	setEventName = func(eventName string) {
		if "" == eventName {
			connection.server.notifier.SendError(
				fmt.Errorf("renaming a server sent event (`%s`) to an empty string is not allowed", connection.eventName),
			)
			return
		}

		connection.eventName = eventName
	}
	return
}

// SendWsUpgrade upgrades to web sockets.
func (connection *Connection) SendWsUpgrade() {
	connection.SendConfiguredWsUpgrade(websocket.Upgrader{
		ReadBufferSize:  10 * KB,
		WriteBufferSize: 10 * KB,
	})
}

// SendConfiguredWsUpgrade upgrades to web sockets.
func (connection *Connection) SendConfiguredWsUpgrade(upgrader websocket.Upgrader) {
	conn, upgradeError := upgrader.Upgrade(connection.writer, connection.request, nil)
	if nil != upgradeError {
		connection.server.notifier.SendErrorAndTrace(upgradeError, 1)
		return
	}
	defer func(conn *websocket.Conn) {
		closeError := conn.Close()
		if nil != closeError {
			connection.server.notifier.SendErrorAndTrace(closeError, 1)
		}
	}(conn)
	connection.webSocket = conn
	connection.locked = true
	return
}

// SendView sends a view.
func (connection *Connection) SendView(view View) {
	if "" != connection.header.Get("Location") {
		return
	}

	if nil == view.Data {
		view.Data = map[string]any{}
	}

	if connection.VerifyAccept("application/json") {
		connection.SendJson(view)
		return
	}

	if "" == view.server {
		view.server = connection.server.viewServer
	}

	if "" == view.index {
		view.index = connection.server.viewIndex
	}

	if "" == view.root {
		view.root = connection.server.viewRoot
	}

	content, compileError := view.Render(connection.server.efs)
	if nil != compileError {
		connection.server.notifier.SendErrorAndTrace(compileError, 1)
		return
	}

	if "" == connection.header.Get("Content-Type") {
		connection.SendHeader("Content-Type", "text/html")
	}

	connection.SendMessage(content)
}
