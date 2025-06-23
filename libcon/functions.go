package libcon

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/libfs"
	"github.com/razshare/frizzante/libglobals"
	"github.com/razshare/frizzante/libmime"
	"github.com/razshare/frizzante/libview"
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
func (connection *Connection) ReceiveCancellation() <-chan struct{} {
	return connection.Request.Context().Done()
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
	cookie, cookieError := connection.Request.Cookie(key)
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
	if connection.WebSocket != nil {
		_, readBytes, readError := connection.WebSocket.ReadMessage()
		if nil != readError {
			connection.Notifier.SendErrorAndTrace(readError, 1)
			return ""
		}
		return string(readBytes)
	}

	readBytes, readAllError := io.ReadAll(connection.Request.Body)
	if nil != readAllError {
		connection.Notifier.SendErrorAndTrace(readAllError, 1)
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
	if connection.WebSocket != nil {
		jsonError := connection.WebSocket.ReadJSON(v)
		if nil != jsonError {
			connection.Notifier.SendErrorAndTrace(jsonError, 1)
			return false
		}
		return true
	}

	readBytes, readAllError := io.ReadAll(connection.Request.Body)
	if nil != readAllError {
		connection.Notifier.SendErrorAndTrace(readAllError, 1)
		return false
	}
	unmarshalError := json.Unmarshal(readBytes, v)
	if nil != unmarshalError {
		connection.Notifier.SendErrorAndTrace(unmarshalError, 1)
		return false
	}
	return true
}

// ReceiveForm reads the message as a form and returns the value.
//
// It silently discards malformed values.
func (connection *Connection) ReceiveForm() url.Values {
	return connection.ReceiveFormWithMaxMemory(2 * libglobals.MB)
}

// ReceiveFormWithMaxMemory reads the message as a form and returns the value.
//
// It silently discards malformed values.
func (connection *Connection) ReceiveFormWithMaxMemory(maxMemory int64) url.Values {
	if connection.WebSocket != nil {
		return url.Values{}
	}

	parseMultipartFormError := connection.Request.ParseMultipartForm(maxMemory)
	if nil != parseMultipartFormError {
		if !errors.Is(parseMultipartFormError, http.ErrNotMultipart) {
			return url.Values{}
		}

		parseFormError := connection.Request.ParseForm()
		if nil != parseFormError {
			return url.Values{}
		}
	}

	return connection.Request.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveQuery(name string) string {
	return connection.Request.URL.Query().Get(name)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceivePath(name string) string {
	return connection.Request.PathValue(name)
}

// ReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveHeader(key string) string {
	return connection.Request.Header.Get(key)
}

// ReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ReceiveContentType(connection *Connection) string {
	return connection.Request.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func VerifyContentType(connection *Connection, contentTypes ...string) bool {
	requestedMime := connection.Request.Header.Get("Content-Type")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func (connection *Connection) VerifyAccept(contentTypes ...string) bool {
	requestedAcceptMime := connection.Request.Header.Get("Accept")
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
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", connection.EventId, connection.EventName)

	_, writeEventError := connection.Writer.Write([]byte(header))
	if nil != writeEventError {
		connection.Notifier.SendErrorAndTrace(writeEventError, 1)
		return
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeEventError = connection.Writer.Write([]byte("data: "))
		if nil != writeEventError {
			connection.Notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}

		_, writeEventError = connection.Writer.Write(line)
		if nil != writeEventError {
			connection.Notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}

		_, writeEventError = connection.Writer.Write([]byte("\r\n"))
		if nil != writeEventError {
			connection.Notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}
	}

	_, writeEventError = connection.Writer.Write([]byte("\r\n"))
	if nil != writeEventError {
		connection.Notifier.SendErrorAndTrace(writeEventError, 1)
		return
	}

	flusher, flushedOk := connection.Writer.(http.Flusher)
	if !flushedOk {
		connection.Notifier.SendErrorAndTrace(errors.New("could not retrieve flusher"), 1)
		return
	}

	flusher.Flush()

	connection.EventId++
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
	if connection.Locked {
		connection.Notifier.SendErrorAndTrace(errors.New("status is locked"), 1)
	}
	connection.Status = code
}

// SendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func (connection *Connection) SendHeader(key string, value string) {
	if connection.Locked {
		connection.Notifier.SendErrorAndTrace(errors.New("header is locked"), 1)
	}

	connection.Header.Set(key, value)
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
	if !connection.Locked {
		connection.Writer.WriteHeader(connection.Status)
		connection.Locked = true
	}

	if connection.WebSocket != nil {
		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			connection.Notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != connection.EventName {
		connection.SendEventContent(content)
		return
	}

	_, writeError := connection.Writer.Write(content)
	if nil != writeError {
		connection.Notifier.SendErrorAndTrace(writeError, 1)
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
		connection.Notifier.SendErrorAndTrace(marshalError, 1)
		return
	}

	if nil == connection.WebSocket {
		contentType := connection.Header.Get("Content-Type")
		if "" == contentType {
			connection.Header.Set("Content-Type", "application/json")
		}
	}

	connection.SendContent(content)
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func (connection *Connection) SendEmbeddedFileOrElse(efs embed.FS, orElse func()) {
	fileName := connection.PublicRoot + connection.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !libfs.EfsIsFile(efs, fileName) || libfs.EfsIsDirectory(efs, fileName) {
		orElse()
		return
	}

	reader, info, readerError := libfs.EfsFileReader(efs, fileName)
	if nil != readerError {
		connection.Notifier.SendErrorAndTrace(readerError, 1)
		return
	}

	if connection.WebSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.Notifier.SendErrorAndTrace(readError, 1)
			return
		}
		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			connection.Notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != connection.EventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.Notifier.SendErrorAndTrace(readError, 1)
		}
		connection.SendEventContent(content)
		return
	}

	if "" == connection.Header.Get("Content-Type") {
		connection.SendHeader("Content-Type", libmime.Mime(fileName))
	}

	if "" == connection.Header.Get("Content-Length") {
		connection.SendHeader("Content-Length", fmt.Sprintf("%d", info.Size()))
	}
	http.ServeContent(connection.Writer, connection.Request, fileName, info.ModTime(), reader)
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func (connection *Connection) SendFileOrElse(orElse func()) {
	fileName := filepath.Join(connection.PublicRoot, connection.Request.RequestURI)

	if !libfs.IsFile(fileName) || libfs.IsDirectory(fileName) {
		connection.SendEmbeddedFileOrElse(connection.Efs, orElse)
		return
	}

	reader, info, readerError := libfs.FileReader(fileName)
	if nil != readerError {
		connection.Notifier.SendErrorAndTrace(readerError, 1)
		return
	}

	if connection.WebSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.Notifier.SendErrorAndTrace(readError, 1)
			return
		}
		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			connection.Notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != connection.EventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			connection.Notifier.SendErrorAndTrace(readError, 1)
			return
		}
		connection.SendEventContent(content)
	}

	if "" == connection.Header.Get("Content-Type") {
		connection.SendHeader("Content-Type", libmime.Mime(fileName))
	}

	if "" == connection.Header.Get("Content-Length") {
		connection.SendHeader("Content-Length", fmt.Sprintf("%d", info.Size()))
	}
	http.ServeContent(connection.Writer, connection.Request, fileName, info.ModTime(), reader)
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
	connection.EventName = "message"
	return func(eventName string) {
		if "" == eventName {
			connection.Notifier.SendErrorAndTrace(
				fmt.Errorf("renaming a server sent event (`%s`) to an empty string is not allowed", connection.EventName),
				1,
			)
		}
		connection.EventName = eventName
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
	connection.EventName = "message"
	setEventName = func(eventName string) {
		if "" == eventName {
			connection.Notifier.SendError(
				fmt.Errorf("renaming a server sent event (`%s`) to an empty string is not allowed", connection.EventName),
			)
			return
		}

		connection.EventName = eventName
	}
	return
}

// SendWsUpgrade upgrades to web sockets.
func (connection *Connection) SendWsUpgrade() {
	connection.SendConfiguredWsUpgrade(websocket.Upgrader{
		ReadBufferSize:  10 * libglobals.KB,
		WriteBufferSize: 10 * libglobals.KB,
	})
}

// SendConfiguredWsUpgrade upgrades to web sockets.
func (connection *Connection) SendConfiguredWsUpgrade(upgrader websocket.Upgrader) {
	conn, upgradeError := upgrader.Upgrade(connection.Writer, connection.Request, nil)
	if nil != upgradeError {
		connection.Notifier.SendErrorAndTrace(upgradeError, 1)
		return
	}
	defer func(conn *websocket.Conn) {
		closeError := conn.Close()
		if nil != closeError {
			connection.Notifier.SendErrorAndTrace(closeError, 1)
		}
	}(conn)
	connection.WebSocket = conn
	connection.Locked = true
	return
}

// SendView sends a view.
func (connection *Connection) SendView(view libview.View) {
	if "" != connection.Header.Get("Location") {
		return
	}

	if nil == view.Data {
		view.Data = map[string]any{}
	}

	if connection.VerifyAccept("application/json") {
		connection.SendJson(map[string]any{
			"name":       view.Name,
			"data":       view.Data,
			"renderMode": view.RenderMode,
		})
		return
	}

	if view.Server == "" {
		view.Server = connection.ViewServer
	}

	if view.Index == "" {
		view.Index = connection.ViewIndex
	}

	if view.Root == "" {
		view.Root = connection.ViewRoot
	}

	content, compileError := view.Render(connection.Efs)
	if nil != compileError {
		connection.Notifier.SendErrorAndTrace(compileError, 1)
		return
	}

	if "" == connection.Header.Get("Content-Type") {
		connection.SendHeader("Content-Type", "text/html")
	}

	connection.SendMessage(content)
}
