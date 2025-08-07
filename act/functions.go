package act

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/mimes"
	"github.com/razshare/frizzante/stack"
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

// ReceiveSessionId tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func ReceiveSessionId(c *connections.Connection) string {
	var id string
	cookies := c.Request.CookiesNamed("session-id")
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
		c.ErrorLog.Println(idObjectError, stack.Trace())
		return ""
	}

	id = idObject.String()

	SendCookie(c, "session-id", id)

	return id
}

// ReceiveCancellation returns a channel that closes when the request gets cancelled.
func ReceiveCancellation(c *connections.Connection) <-chan struct{} {
	return c.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(c *connections.Connection) *bool {
	isAlive := true
	go func() {
		<-ReceiveCancellation(c)
		isAlive = false
	}()
	return &isAlive
}

// ReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func ReceiveCookie(c *connections.Connection, key string) string {
	cookie, cookieError := c.Request.Cookie(key)
	if cookieError != nil {
		c.ErrorLog.Println(cookieError, stack.Trace())
		return ""
	}

	data, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		c.ErrorLog.Println(queryError, stack.Trace())
		return ""
	}

	return data
}

// ReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func ReceiveMessage(c *connections.Connection) string {
	if c.WebSocket != nil {
		_, data, readError := c.WebSocket.ReadMessage()
		if readError != nil {
			c.ErrorLog.Println(readError, stack.Trace())
			return ""
		}
		return string(data)
	}

	data, readError := io.ReadAll(c.Request.Body)
	if readError != nil {
		c.ErrorLog.Println(readError, stack.Trace())
		return ""
	}
	return string(data)
}

// ReceiveJson reads the next JSON-encoded message from the
// c and stores it in the value pointed to by va.
//
// Compatible with web sockets.
func ReceiveJson(c *connections.Connection, v any) {
	if c.WebSocket != nil {
		jsonError := c.WebSocket.ReadJSON(v)
		if jsonError != nil {
			c.ErrorLog.Println(jsonError, stack.Trace())
			return
		}
		return
	}

	data, readError := io.ReadAll(c.Request.Body)
	if readError != nil {
		c.ErrorLog.Println(readError, stack.Trace())
		return
	}

	jsonError := json.Unmarshal(data, v)
	if jsonError != nil {
		c.ErrorLog.Println(jsonError, stack.Trace())
		return
	}
}

// ReceiveForm reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of 2MB
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func ReceiveForm(c *connections.Connection) url.Values {
	return ReceiveFormWithMaxMemory(c, 2*globals.MB)
}

// ReceiveFormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of maxMemory bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func ReceiveFormWithMaxMemory(c *connections.Connection, m int64) url.Values {
	if c.WebSocket != nil {
		c.ErrorLog.Println(errors.New("c is not of type web socket"), stack.Trace())
		return url.Values{}
	}

	formError := c.Request.ParseMultipartForm(m)
	if formError != nil {
		if !errors.Is(formError, http.ErrNotMultipart) {
			return url.Values{}
		}

		formError = c.Request.ParseForm()
		if formError != nil {
			c.ErrorLog.Println(formError, stack.Trace())
			return url.Values{}
		}
	}

	return c.Request.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func ReceiveQuery(c *connections.Connection, k string) string {
	return c.Request.URL.Query().Get(k)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func ReceivePath(c *connections.Connection, k string) string {
	return c.Request.PathValue(k)
}

// ReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func ReceiveHeader(c *connections.Connection, k string) string {
	return c.Request.Header.Get(k)
}

// ReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ReceiveContentType(c *connections.Connection) string {
	return c.Request.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func VerifyContentType(c *connections.Connection, t ...string) bool {
	requestedMime := c.Request.Header.Get("Content-Type")
	for _, acceptedMime := range t {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func VerifyAccept(c *connections.Connection, t ...string) bool {
	requestedAcceptMime := c.Request.Header.Get("Accept")
	for _, acceptedMime := range t {
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
func SendEventContent(c *connections.Connection, d []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", c.EventId, c.EventName)

	_, writeError := c.Writer.Write([]byte(header))
	if writeError != nil {
		c.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	for _, line := range bytes.Split(d, []byte("\r\n")) {
		_, writeError = c.Writer.Write([]byte("data: "))
		if writeError != nil {
			c.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = c.Writer.Write(line)
		if writeError != nil {
			c.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = c.Writer.Write([]byte("\r\n"))
		if writeError != nil {
			c.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	_, writeError = c.Writer.Write([]byte("\r\n"))
	if writeError != nil {
		c.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	flusher, flushedOk := c.Writer.(http.Flusher)
	if !flushedOk {
		c.ErrorLog.Println(errors.New("could not retrieve flusher"), stack.Trace())
		return
	}

	flusher.Flush()

	c.EventId++
}

// SendNavigate redirects the request to a location with status 302.
func SendNavigate(c *connections.Connection, l string) {
	SendRedirect(c, l, 302)
	SendFlush(c)
}

// SendRedirect redirects the request to a location with a status.
func SendRedirect(c *connections.Connection, l string, status int) {
	SendStatus(c, status)
	SendHeader(c, "Location", l)
}

// SendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func SendStatus(c *connections.Connection, s int) {
	if c.Locked {
		c.ErrorLog.Println("status is locked", stack.Trace())
		return
	}

	c.Status = s
}

// SendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func SendHeader(c *connections.Connection, k string, v string) {
	if c.Locked {
		c.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	c.Writer.Header().Set(k, v)
}

func SendHeaders(c *connections.Connection, h map[string]string) {
	if c.Locked {
		c.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	for key, value := range h {
		c.Writer.Header().Set(key, value)
	}
}

// SendContentType sets the Content-Type header field.
func SendContentType(c *connections.Connection, t string) {
	SendHeader(c, "Content-Type", t)
}

// SendCookie sends a cookies to the client.
func SendCookie(c *connections.Connection, key string, value string) {
	SendHeader(c, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}

func SendFlush(c *connections.Connection) {
	SendMessage(c, "")
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
func SendContent(c *connections.Connection, d []byte) {
	if !c.Locked {
		c.Writer.WriteHeader(c.Status)
		c.Locked = true
	}

	if c.WebSocket != nil {
		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, d)
		if writeError != nil {
			c.ErrorLog.Println(writeError, stack.Trace())
		}
		return
	}

	if "" != c.EventName {
		SendEventContent(c, d)
		return
	}

	_, writeError := c.Writer.Write(d)
	if writeError != nil {
		c.ErrorLog.Println(writeError, stack.Trace())
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
func SendMessage(c *connections.Connection, m string) {
	SendContent(c, []byte(m))
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
func SendMessagef(c *connections.Connection, f string, v ...any) {
	SendContent(c, []byte(fmt.Sprintf(f, v...)))
}

// SendNotFound sends a message with status 404 Not Found.
func SendNotFound(c *connections.Connection, m string) {
	SendStatus(c, http.StatusNotFound)
	SendMessage(c, m)
}

// SendUnauthorized sends a message with status 401 Unauthorized.
func SendUnauthorized(c *connections.Connection, m string) {
	SendStatus(c, http.StatusUnauthorized)
	SendMessage(c, m)
}

// SendBadRequest sends a message with status 400 Bad Request.
func SendBadRequest(c *connections.Connection, m string) {
	SendStatus(c, http.StatusBadRequest)
	SendMessage(c, m)
}

// SendError sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func SendError(c *connections.Connection, e error) {
	SendStatus(c, http.StatusBadRequest)
	SendMessage(c, e.Error())
}

// SendForbidden sends a message with status 403 Forbidden.
func SendForbidden(c *connections.Connection, m string) {
	SendStatus(c, http.StatusForbidden)
	SendMessage(c, m)
}

// SendTooManyRequests sends a message with status 403 Forbidden.
func SendTooManyRequests(c *connections.Connection, m string) {
	SendStatus(c, http.StatusTooManyRequests)
	SendMessage(c, m)
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
func SendJson(c *connections.Connection, v any) {
	data, jsonError := json.Marshal(v)
	if jsonError != nil {
		c.ErrorLog.Println(jsonError, stack.Trace())
		return
	}

	if nil == c.WebSocket {
		contentType := c.Writer.Header().Get("Content-Type")
		if "" == contentType {
			c.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	SendContent(c, data)
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func SendEmbeddedFileOrElse(c *connections.Connection, efs embed.FS, or func()) {
	fileName := c.PublicRoot + c.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !embeds.IsFile(efs, fileName) || embeds.IsDirectory(efs, fileName) {
		or()
		return
	}

	reader, readerInfo, readerError := embeds.FileReader(efs, fileName)
	if readerError != nil {
		c.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.ErrorLog.Println(writeError, stack.Trace())
			return
		}
		return
	}

	if "" != c.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.ErrorLog.Println(readError, stack.Trace())
			return
		}

		SendEventContent(c, data)
		return
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		SendHeader(c, "Content-Type", mimes.Mime(fileName))
	}

	if "" == c.Writer.Header().Get("Content-Length") {
		SendHeader(c, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(c.Writer, c.Request, fileName, readerInfo.ModTime(), reader)
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func SendFileOrElse(c *connections.Connection, or func()) {
	fileName := filepath.Join(c.PublicRoot, c.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		SendEmbeddedFileOrElse(c, c.Efs, or)
		return
	}

	reader, readerInfo, readerError := files.FileReader(fileName)
	if readerError != nil {
		c.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	if "" != c.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.ErrorLog.Println(readError, stack.Trace())
			return
		}

		SendEventContent(c, data)
		return
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		SendHeader(c, "Content-Type", mimes.Mime(fileName))
	}

	if "" == c.Writer.Header().Get("Content-Length") {
		SendHeader(c, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(c.Writer, c.Request, fileName, readerInfo.ModTime(), reader)
}

// SendSseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SendSseUpgrade(c *connections.Connection) func(string) {
	SendHeaders(c, map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Connection":                    "keep-alive",
	})

	c.EventName = "message"

	return func(eventName string) { c.EventName = eventName }
}

// SendWsUpgrade upgrades to web sockets.
func SendWsUpgrade(c *connections.Connection) {
	SendConfiguredWsUpgrade(c, websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// SendConfiguredWsUpgrade upgrades to web sockets.
func SendConfiguredWsUpgrade(c *connections.Connection, u websocket.Upgrader) {
	webSocketConnection, upgradeError := u.Upgrade(c.Writer, c.Request, nil)
	if upgradeError != nil {
		c.ErrorLog.Println(upgradeError, stack.Trace())
		return
	}

	defer func(webSocketConnection *websocket.Conn) {
		closeError := c.WebSocket.Close()
		if closeError != nil {
			c.ErrorLog.Println(closeError, stack.Trace())
		}
	}(webSocketConnection)

	c.WebSocket = webSocketConnection
	c.Locked = true

	return
}

// SendView sends a view.
func SendView(c *connections.Connection, v views.View) {
	if c.Writer.Header().Get("Location") != "" {
		return
	}

	if VerifyAccept(c, "application/json") {
		if v.Data == nil {
			v.Data = map[string]any{}
		}
		props := map[string]any{
			"name":       v.Name,
			"data":       v.Data,
			"renderMode": v.RenderMode,
		}
		SendJson(c, props)
		return
	}

	html, renderError := views.Render(&v, c.App, c.AppConfig)
	if renderError != nil {
		c.ErrorLog.Println(renderError, stack.Trace())
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		SendHeader(c, "Content-Type", "text/html")
	}

	SendMessage(c, html)
}
