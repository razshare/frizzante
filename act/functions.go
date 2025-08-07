package act

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
	ef "github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/mime"
	"github.com/razshare/frizzante/server"
	"github.com/razshare/frizzante/stack"
	v "github.com/razshare/frizzante/view"
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

// ReceiveSessionId tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func ReceiveSessionId(connection *server.Connection) string {
	var id string
	cookies := connection.Request.CookiesNamed("session-id")
	cookiesCount := 0

	for _, cookie := range cookies {
		id = cookie.Value
		cookiesCount++
	}

	if cookiesCount > 0 {
		return id
	}

	// Create new session.
	idObject, idObjectError := uuid.NewV4()
	if idObjectError != nil {
		connection.Server.ErrorLog.Println(idObjectError, stack.Trace())
		return ""
	}

	id = idObject.String()

	SendCookie(connection, "session-id", id)

	return id
}

// ReceiveCancellation returns a channel that closes when the request gets cancelled.
func ReceiveCancellation(connection *server.Connection) <-chan struct{} {
	return connection.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(connection *server.Connection) *bool {
	isAlive := true
	go func() {
		<-ReceiveCancellation(connection)
		isAlive = false
	}()
	return &isAlive
}

// ReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func ReceiveCookie(connection *server.Connection, key string) string {
	cookie, cookieError := connection.Request.Cookie(key)
	if cookieError != nil {
		connection.Server.ErrorLog.Println(cookieError, stack.Trace())
		return ""
	}

	data, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		connection.Server.ErrorLog.Println(queryError, stack.Trace())
		return ""
	}

	return data
}

// ReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func ReceiveMessage(connection *server.Connection) string {
	if connection.WebSocket != nil {
		_, data, readError := connection.WebSocket.ReadMessage()
		if readError != nil {
			connection.Server.ErrorLog.Println(readError, stack.Trace())
			return ""
		}
		return string(data)
	}

	data, readError := io.ReadAll(connection.Request.Body)
	if readError != nil {
		connection.Server.ErrorLog.Println(readError, stack.Trace())
		return ""
	}
	return string(data)
}

// ReceiveJson reads the next JSON-encoded message from the
// connection and stores it in the value pointed to by val.
//
// Compatible with web sockets.
func ReceiveJson(connection *server.Connection, value any) {
	if connection.WebSocket != nil {
		jsonError := connection.WebSocket.ReadJSON(value)
		if jsonError != nil {
			connection.Server.ErrorLog.Println(jsonError, stack.Trace())
			return
		}
		return
	}

	data, readError := io.ReadAll(connection.Request.Body)
	if readError != nil {
		connection.Server.ErrorLog.Println(readError, stack.Trace())
		return
	}

	jsonError := json.Unmarshal(data, value)
	if jsonError != nil {
		connection.Server.ErrorLog.Println(jsonError, stack.Trace())
		return
	}
}

// ReceiveForm reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of 2MB
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func ReceiveForm(connection *server.Connection) url.Values {
	return ReceiveFormWithMaxMemory(connection, 2*globals.MB)
}

// ReceiveFormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of maxMemory bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func ReceiveFormWithMaxMemory(connection *server.Connection, maxMemory int64) url.Values {
	if connection.WebSocket != nil {
		connection.Server.ErrorLog.Println(errors.New("connection is not of type web socket"), stack.Trace())
		return url.Values{}
	}

	formError := connection.Request.ParseMultipartForm(maxMemory)
	if formError != nil {
		if !errors.Is(formError, http.ErrNotMultipart) {
			return url.Values{}
		}

		formError = connection.Request.ParseForm()
		if formError != nil {
			connection.Server.ErrorLog.Println(formError, stack.Trace())
			return url.Values{}
		}
	}

	return connection.Request.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func ReceiveQuery(connection *server.Connection, key string) string {
	return connection.Request.URL.Query().Get(key)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func ReceivePath(connection *server.Connection, key string) string {
	return connection.Request.PathValue(key)
}

// ReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func ReceiveHeader(connection *server.Connection, key string) string {
	return connection.Request.Header.Get(key)
}

// ReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ReceiveContentType(connection *server.Connection) string {
	return connection.Request.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func VerifyContentType(connection *server.Connection, contentType ...string) bool {
	requestedMime := connection.Request.Header.Get("Content-Type")
	for _, acceptedMime := range contentType {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func VerifyAccept(connection *server.Connection, accept ...string) bool {
	requestedAcceptMime := connection.Request.Header.Get("Accept")
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
func SendEventContent(connection *server.Connection, content []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", connection.EventId, connection.EventName)

	_, writeError := connection.Writer.Write([]byte(header))
	if writeError != nil {
		connection.Server.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeError = connection.Writer.Write([]byte("data: "))
		if writeError != nil {
			connection.Server.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = connection.Writer.Write(line)
		if writeError != nil {
			connection.Server.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = connection.Writer.Write([]byte("\r\n"))
		if writeError != nil {
			connection.Server.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	_, writeError = connection.Writer.Write([]byte("\r\n"))
	if writeError != nil {
		connection.Server.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	flusher, flushedOk := connection.Writer.(http.Flusher)
	if !flushedOk {
		connection.Server.ErrorLog.Println(errors.New("could not retrieve flusher"), stack.Trace())
		return
	}

	flusher.Flush()

	connection.EventId++
}

// SendNavigate redirects the request to a location with status 302.
func SendNavigate(connection *server.Connection, location string) {
	SendRedirect(connection, location, 302)
	SendFlush(connection)
}

// SendRedirect redirects the request to a location with a status.
func SendRedirect(connection *server.Connection, location string, status int) {
	SendStatus(connection, status)
	SendHeader(connection, "Location", location)
}

// SendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func SendStatus(connection *server.Connection, status int) {
	if connection.Locked {
		connection.Server.ErrorLog.Println("status is locked", stack.Trace())
		return
	}

	connection.Status = status
}

// SendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func SendHeader(connection *server.Connection, key string, value string) {
	if connection.Locked {
		connection.Server.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	connection.Writer.Header().Set(key, value)
}

func SendHeaders(connection *server.Connection, headers map[string]string) {
	if connection.Locked {
		connection.Server.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	for key, value := range headers {
		connection.Writer.Header().Set(key, value)
	}
}

// SendContentType sets the Content-Type header field.
func SendContentType(connection *server.Connection, contentType string) {
	SendHeader(connection, "Content-Type", contentType)
}

// SendCookie sends a cookies to the client.
func SendCookie(connection *server.Connection, key string, value string) {
	SendHeader(connection, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}

func SendFlush(connection *server.Connection) {
	SendMessage(connection, "")
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
func SendContent(connection *server.Connection, content []byte) {
	if !connection.Locked {
		connection.Writer.WriteHeader(connection.Status)
		connection.Locked = true
	}

	if connection.WebSocket != nil {
		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			connection.Server.ErrorLog.Println(writeError, stack.Trace())
		}
		return
	}

	if "" != connection.EventName {
		SendEventContent(connection, content)
		return
	}

	_, writeError := connection.Writer.Write(content)
	if writeError != nil {
		connection.Server.ErrorLog.Println(writeError, stack.Trace())
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
func SendMessage(connection *server.Connection, message string) {
	SendContent(connection, []byte(message))
}

// SendMessagef sends utf-8 safe content using a format.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func SendMessagef(connection *server.Connection, format string, values ...any) {
	SendContent(connection, []byte(fmt.Sprintf(format, values...)))
}

// SendNotFound sends a message with status 404 Not Found.
func SendNotFound(connection *server.Connection, message string) {
	SendStatus(connection, http.StatusNotFound)
	SendMessage(connection, message)
}

// SendUnauthorized sends a message with status 401 Unauthorized.
func SendUnauthorized(connection *server.Connection, message string) {
	SendStatus(connection, http.StatusUnauthorized)
	SendMessage(connection, message)
}

// SendBadRequest sends a message with status 400 Bad Request.
func SendBadRequest(connection *server.Connection, message string) {
	SendStatus(connection, http.StatusBadRequest)
	SendMessage(connection, message)
}

// SendError sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func SendError(connection *server.Connection, err error) {
	SendStatus(connection, http.StatusBadRequest)
	SendMessage(connection, err.Error())
}

// SendForbidden sends a message with status 403 Forbidden.
func SendForbidden(connection *server.Connection, message string) {
	SendStatus(connection, http.StatusForbidden)
	SendMessage(connection, message)
}

// SendTooManyRequests sends a message with status 403 Forbidden.
func SendTooManyRequests(connection *server.Connection, message string) {
	SendStatus(connection, http.StatusTooManyRequests)
	SendMessage(connection, message)
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
func SendJson(connection *server.Connection, value any) {
	data, jsonError := json.Marshal(value)
	if jsonError != nil {
		connection.Server.ErrorLog.Println(jsonError, stack.Trace())
		return
	}

	if nil == connection.WebSocket {
		contentType := connection.Writer.Header().Get("Content-Type")
		if "" == contentType {
			connection.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	SendContent(connection, data)
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func SendEmbeddedFileOrElse(connection *server.Connection, efs embed.FS, orElse func()) {
	fileName := connection.Server.PublicRoot + connection.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !ef.IsFile(efs, fileName) || ef.IsDirectory(efs, fileName) {
		orElse()
		return
	}

	reader, readerInfo, readerError := ef.FileReader(efs, fileName)
	if readerError != nil {
		connection.Server.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if connection.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			connection.Server.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			connection.Server.ErrorLog.Println(writeError, stack.Trace())
			return
		}
		return
	}

	if "" != connection.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			connection.Server.ErrorLog.Println(readError, stack.Trace())
			return
		}

		SendEventContent(connection, data)
		return
	}

	if "" == connection.Writer.Header().Get("Content-Type") {
		SendHeader(connection, "Content-Type", mime.Mime(fileName))
	}

	if "" == connection.Writer.Header().Get("Content-Length") {
		SendHeader(connection, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(connection.Writer, connection.Request, fileName, readerInfo.ModTime(), reader)
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func SendFileOrElse(connection *server.Connection, orElse func()) {
	fileName := filepath.Join(connection.Server.PublicRoot, connection.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		SendEmbeddedFileOrElse(connection, connection.Server.Efs, orElse)
		return
	}

	reader, readerInfo, readerError := files.FileReader(fileName)
	if readerError != nil {
		connection.Server.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if connection.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			connection.Server.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			connection.Server.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	if "" != connection.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			connection.Server.ErrorLog.Println(readError, stack.Trace())
			return
		}

		SendEventContent(connection, data)
		return
	}

	if "" == connection.Writer.Header().Get("Content-Type") {
		SendHeader(connection, "Content-Type", mime.Mime(fileName))
	}

	if "" == connection.Writer.Header().Get("Content-Length") {
		SendHeader(connection, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(connection.Writer, connection.Request, fileName, readerInfo.ModTime(), reader)
}

// SendSseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SendSseUpgrade(connection *server.Connection) func(eventName string) {
	SendHeaders(connection, map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Connection":                    "keep-alive",
	})

	connection.EventName = "message"

	return func(eventName string) { connection.EventName = eventName }
}

// SendWsUpgrade upgrades to web sockets.
func SendWsUpgrade(connection *server.Connection) {
	SendConfiguredWsUpgrade(connection, websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// SendConfiguredWsUpgrade upgrades to web sockets.
func SendConfiguredWsUpgrade(connection *server.Connection, upgrader websocket.Upgrader) {
	webSocketConnection, upgradeError := upgrader.Upgrade(connection.Writer, connection.Request, nil)
	if upgradeError != nil {
		connection.Server.ErrorLog.Println(upgradeError, stack.Trace())
		return
	}

	defer func(webSocketConnection *websocket.Conn) {
		closeError := connection.WebSocket.Close()
		if closeError != nil {
			connection.Server.ErrorLog.Println(closeError, stack.Trace())
		}
	}(webSocketConnection)

	connection.WebSocket = webSocketConnection
	connection.Locked = true

	return
}

// SendView sends a view.
func SendView(connection *server.Connection, view v.View) {
	if connection.Writer.Header().Get("Location") != "" {
		return
	}

	if VerifyAccept(connection, "application/json") {
		if view.Data == nil {
			view.Data = map[string]any{}
		}
		props := map[string]any{
			"name":       view.Name,
			"data":       view.Data,
			"renderMode": view.RenderMode,
		}
		SendJson(connection, props)
		return
	}

	html, renderError := v.Render(&view, connection.Application, connection.Server.Application)
	if renderError != nil {
		connection.Server.ErrorLog.Println(renderError, stack.Trace())
	}

	if "" == connection.Writer.Header().Get("Content-Type") {
		SendHeader(connection, "Content-Type", "text/html")
	}

	SendMessage(connection, html)
}
