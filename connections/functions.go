package connections

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/mimes"
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
func (con *Connection) ReceiveCancellation() <-chan struct{} {
	return con.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func (con *Connection) IsAlive() *bool {
	value := true
	go func() {
		<-con.ReceiveCancellation()
		value = false
	}()
	return &value
}

// ReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func (con *Connection) ReceiveCookie(key string) string {
	cookie, cookieError := con.Request.Cookie(key)
	if cookieError != nil {
		con.Notifier.SendErrorAndTrace(cookieError, 1)
		return ""
	}
	value, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		con.Notifier.SendErrorAndTrace(queryError, 1)
		return ""
	}

	return value
}

// ReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func (con *Connection) ReceiveMessage() string {
	if con.WebSocket != nil {
		_, readBytes, readError := con.WebSocket.ReadMessage()
		if readError != nil {
			con.Notifier.SendErrorAndTrace(readError, 1)
			return ""
		}
		return string(readBytes)
	}

	readBytes, readAllError := io.ReadAll(con.Request.Body)
	if readAllError != nil {
		con.Notifier.SendErrorAndTrace(readAllError, 1)
		return ""
	}
	return string(readBytes)
}

// ReceiveJson reads the next JSON-encoded message from the
// connection and stores it in the value pointed to by val.
//
// Compatible with web sockets.
func (con *Connection) ReceiveJson(val any) error {
	if con.WebSocket != nil {
		jsonError := con.WebSocket.ReadJSON(val)
		if jsonError != nil {
			return jsonError
		}
		return nil
	}

	readBytes, readAllError := io.ReadAll(con.Request.Body)
	if readAllError != nil {
		return readAllError
	}
	unmarshalError := json.Unmarshal(readBytes, val)
	if unmarshalError != nil {
		return unmarshalError
	}
	return nil
}

// ReceiveForm reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of 2MB
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func (con *Connection) ReceiveForm() url.Values {
	return con.ReceiveFormWithMaxMemory(2 * globals.MB)
}

// ReceiveFormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of maxMemory bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func (con *Connection) ReceiveFormWithMaxMemory(maxMemory int64) url.Values {
	if con.WebSocket != nil {
		return url.Values{}
	}

	parseMultipartFormError := con.Request.ParseMultipartForm(maxMemory)
	if parseMultipartFormError != nil {
		if !errors.Is(parseMultipartFormError, http.ErrNotMultipart) {
			return url.Values{}
		}

		parseFormError := con.Request.ParseForm()
		if parseFormError != nil {
			con.Notifier.SendErrorAndTrace(parseFormError, 1)
			return url.Values{}
		}
	}

	return con.Request.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func (con *Connection) ReceiveQuery(name string) string {
	return con.Request.URL.Query().Get(name)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func (con *Connection) ReceivePath(name string) string {
	return con.Request.PathValue(name)
}

// ReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func (con *Connection) ReceiveHeader(key string) string {
	return con.Request.Header.Get(key)
}

// ReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func (con *Connection) ReceiveContentType() string {
	return con.Request.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func (con *Connection) VerifyContentType(t ...string) bool {
	requestedMime := con.Request.Header.Get("Content-Type")
	for _, acceptedMime := range t {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func (con *Connection) VerifyAccept(t ...string) bool {
	requestedAcceptMime := con.Request.Header.Get("Accept")
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
func (con *Connection) SendEventContent(val []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", con.EventId, con.EventName)

	_, writeEventError := con.Writer.Write([]byte(header))
	if writeEventError != nil {
		con.Notifier.SendErrorAndTrace(writeEventError, 1)
		return
	}

	for _, line := range bytes.Split(val, []byte("\r\n")) {
		_, writeEventError = con.Writer.Write([]byte("data: "))
		if writeEventError != nil {
			con.Notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}

		_, writeEventError = con.Writer.Write(line)
		if writeEventError != nil {
			con.Notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}

		_, writeEventError = con.Writer.Write([]byte("\r\n"))
		if writeEventError != nil {
			con.Notifier.SendErrorAndTrace(writeEventError, 1)
			return
		}
	}

	_, writeEventError = con.Writer.Write([]byte("\r\n"))
	if writeEventError != nil {
		con.Notifier.SendErrorAndTrace(writeEventError, 1)
		return
	}

	flusher, flushedOk := con.Writer.(http.Flusher)
	if !flushedOk {
		con.Notifier.SendErrorAndTrace(errors.New("could not retrieve flusher"), 1)
		return
	}

	flusher.Flush()

	con.EventId++
}

// SendNavigate redirects the request with status 302.
func (con *Connection) SendNavigate(loc string) {
	con.SendRedirect(loc, 302)
	con.SendFlush()
}

// SendRedirect redirects the request.
func (con *Connection) SendRedirect(loc string, code int) {
	con.SendStatus(code)
	con.SendHeader("Location", loc)
}

// SendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func (con *Connection) SendStatus(code int) {
	if con.Locked {
		con.Notifier.SendErrorAndTrace(errors.New("status is locked"), 1)
	}
	con.Status = code
}

// SendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func (con *Connection) SendHeader(key string, val string) {
	if con.Locked {
		con.Notifier.SendErrorAndTrace(errors.New("header is locked"), 1)
	}

	con.Header.Set(key, val)
}

// SendContentType sets the Content-Type header field.
func (con *Connection) SendContentType(t string) {
	con.SendHeader("Content-Type", t)
}

// SendCookie sends a cookies to the client.
func (con *Connection) SendCookie(key string, val string) {
	con.SendHeader(
		"Set-Cookie",
		fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(val)),
	)
}

func (con *Connection) SendFlush() {
	con.SendMessage("")
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
func (con *Connection) SendContent(val []byte) {
	if !con.Locked {
		con.Writer.WriteHeader(con.Status)
		con.Locked = true
	}

	if con.WebSocket != nil {
		writeError := con.WebSocket.WriteMessage(websocket.TextMessage, val)
		if writeError != nil {
			con.Notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != con.EventName {
		con.SendEventContent(val)
		return
	}

	_, writeError := con.Writer.Write(val)
	if writeError != nil {
		con.Notifier.SendErrorAndTrace(writeError, 1)
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
func (con *Connection) SendMessage(val string) {
	con.SendContent([]byte(val))
}

// SendNotFound sends a message with status 404 Not Found.
func (con *Connection) SendNotFound(val string) {
	con.SendStatus(http.StatusNotFound)
	con.SendMessage(val)
}

// SendUnauthorized sends a message with status 401 Unauthorized.
func (con *Connection) SendUnauthorized(val string) {
	con.SendStatus(http.StatusUnauthorized)
	con.SendMessage(val)
}

// SendBadRequest sends a message with status 400 Bad Request.
func (con *Connection) SendBadRequest(val string) {
	con.SendStatus(http.StatusBadRequest)
	con.SendMessage(val)
}

// SendInternalServerError sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func (con *Connection) SendInternalServerError(err error) {
	con.SendStatus(http.StatusBadRequest)
	con.SendMessage(err.Error())
}

// SendForbidden sends a message with status 403 Forbidden.
func (con *Connection) SendForbidden(val string) {
	con.SendStatus(http.StatusForbidden)
	con.SendMessage(val)
}

// SendTooManyRequests sends a message with status 403 Forbidden.
func (con *Connection) SendTooManyRequests(val string) {
	con.SendStatus(http.StatusTooManyRequests)
	con.SendMessage(val)
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
func (con *Connection) SendJson(val any) {
	content, marshalError := json.Marshal(val)
	if marshalError != nil {
		con.Notifier.SendErrorAndTrace(marshalError, 1)
		return
	}

	if nil == con.WebSocket {
		contentType := con.Header.Get("Content-Type")
		if "" == contentType {
			con.Header.Set("Content-Type", "application/json")
		}
	}

	con.SendContent(content)
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func (con *Connection) SendEmbeddedFileOrElse(efs embed.FS, fun func()) {
	fileName := con.PublicRoot + con.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !embeds.IsFile(efs, fileName) || embeds.IsDirectory(efs, fileName) {
		fun()
		return
	}

	reader, info, readerError := embeds.FileReader(efs, fileName)
	if readerError != nil {
		con.Notifier.SendErrorAndTrace(readerError, 1)
		return
	}

	if con.WebSocket != nil {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			con.Notifier.SendErrorAndTrace(readError, 1)
			return
		}
		writeError := con.WebSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			con.Notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != con.EventName {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			con.Notifier.SendErrorAndTrace(readError, 1)
		}
		con.SendEventContent(content)
		return
	}

	if "" == con.Header.Get("Content-Type") {
		con.SendHeader("Content-Type", mimes.Mime(fileName))
	}

	if "" == con.Header.Get("Content-Length") {
		con.SendHeader("Content-Length", fmt.Sprintf("%d", info.Size()))
	}
	http.ServeContent(con.Writer, con.Request, fileName, info.ModTime(), reader)
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func (con *Connection) SendFileOrElse(orElse func()) {
	fileName := filepath.Join(con.PublicRoot, con.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		con.SendEmbeddedFileOrElse(con.Efs, orElse)
		return
	}

	reader, info, readerError := files.FileReader(fileName)
	if readerError != nil {
		con.Notifier.SendErrorAndTrace(readerError, 1)
		return
	}

	if con.WebSocket != nil {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			con.Notifier.SendErrorAndTrace(readError, 1)
			return
		}
		writeError := con.WebSocket.WriteMessage(websocket.TextMessage, content)
		if writeError != nil {
			con.Notifier.SendErrorAndTrace(writeError, 1)
		}
		return
	}

	if "" != con.EventName {
		content, readError := io.ReadAll(reader)
		if readError != nil {
			con.Notifier.SendErrorAndTrace(readError, 1)
			return
		}
		con.SendEventContent(content)
	}

	if "" == con.Header.Get("Content-Type") {
		con.SendHeader("Content-Type", mimes.Mime(fileName))
	}

	if "" == con.Header.Get("Content-Length") {
		con.SendHeader("Content-Length", fmt.Sprintf("%d", info.Size()))
	}
	http.ServeContent(con.Writer, con.Request, fileName, info.ModTime(), reader)
}

// SendSseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func (con *Connection) SendSseUpgrade() (ev func(name string)) {
	con.SendHeader("Access-Control-Allow-Origin", "*")
	con.SendHeader("Access-Control-Expose-Headers", "Content-Type")
	con.SendHeader("Content-Type", "text/event-stream")
	con.SendHeader("Cache-Control", "no-cache")
	con.SendHeader("Connection", "keep-alive")
	con.EventName = "message"
	ev = func(name string) {
		if "" == name {
			con.Notifier.SendError(
				fmt.Errorf("renaming a server sent event (`%s`) to an empty string is not allowed", con.EventName),
			)
			return
		}

		con.EventName = name
	}
	return
}

// SendWsUpgrade upgrades to web sockets.
func (con *Connection) SendWsUpgrade() {
	con.SendConfiguredWsUpgrade(websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// SendConfiguredWsUpgrade upgrades to web sockets.
func (con *Connection) SendConfiguredWsUpgrade(up websocket.Upgrader) {
	conn, upgradeError := up.Upgrade(con.Writer, con.Request, nil)
	if upgradeError != nil {
		con.Notifier.SendErrorAndTrace(upgradeError, 1)
		return
	}
	defer func(conn *websocket.Conn) {
		closeError := conn.Close()
		if closeError != nil {
			con.Notifier.SendErrorAndTrace(closeError, 1)
		}
	}(conn)
	con.WebSocket = conn
	con.Locked = true
	return
}

// SendView sends a view.
func (con *Connection) SendView(val views.View) {
	if con.Header.Get("Location") != "" {
		return
	}

	if con.VerifyAccept("application/json") {
		if val.Data == nil {
			val.Data = map[string]any{}
		}
		props := map[string]any{
			"name":       val.Name,
			"data":       val.Data,
			"renderMode": val.RenderMode,
		}
		con.SendJson(props)
		return
	}

	if val.ServerJs == "" {
		val.ServerJs = con.ServerJs
	}

	if val.IndexHtml == "" {
		val.IndexHtml = con.IndexHtml
	}

	if val.AppRoot == "" {
		val.AppRoot = con.AppRoot
	}

	txt, renderError := val.Render(con.Efs)
	if renderError != nil {
		con.Notifier.SendErrorAndTrace(renderError, 1)
		return
	}

	if "" == con.Header.Get("Content-Type") {
		con.SendHeader("Content-Type", "text/html")
	}

	con.SendMessage(txt)
}
