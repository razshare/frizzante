package send

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/mime"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/view"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// EventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a Server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func EventContent(c *conn.Conn, d []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", c.EventId, c.EventName)

	_, writeError := c.Writer.Write([]byte(header))
	if writeError != nil {
		c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	for _, line := range bytes.Split(d, []byte("\r\n")) {
		_, writeError = c.Writer.Write([]byte("data: "))
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = c.Writer.Write(line)
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = c.Writer.Write([]byte("\r\n"))
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	_, writeError = c.Writer.Write([]byte("\r\n"))
	if writeError != nil {
		c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	flusher, flushedOk := c.Writer.(http.Flusher)
	if !flushedOk {
		c.Container.Config.ErrorLog.Println(errors.New("could not retrieve flusher"), stack.Trace())
		return
	}

	flusher.Flush()

	c.EventId++
}

// Navigate redirects the request to a location with status 302.
func Navigate(c *conn.Conn, l string) {
	Redirect(c, l, 302)
	Flush(c)
}

// Redirect redirects the request to a location with a status.
func Redirect(c *conn.Conn, l string, status int) {
	Status(c, status)
	Header(c, "Location", l)
}

// Status sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func Status(c *conn.Conn, s int) {
	if c.Locked {
		c.Container.Config.ErrorLog.Println("status is locked", stack.Trace())
		return
	}

	c.Status = s
}

// Header sends a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func Header(c *conn.Conn, k string, v string) {
	if c.Locked {
		c.Container.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	c.Writer.Header().Set(k, v)
}

// Headers sends header fields.
func Headers(c *conn.Conn, h map[string]string) {
	if c.Locked {
		c.Container.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	for key, value := range h {
		c.Writer.Header().Set(key, value)
	}
}

// ContentType sets the Content-Type header field.
func ContentType(c *conn.Conn, t string) {
	Header(c, "Content-Type", t)
}

// Cookie sends a cookies to the client.
func Cookie(c *conn.Conn, key string, value string) {
	Header(c, "Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
}

func Flush(c *conn.Conn) {
	Message(c, "")
}

// Content sends binary safe content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func Content(c *conn.Conn, d []byte) {
	if !c.Locked {
		c.Writer.WriteHeader(c.Status)
		c.Locked = true
	}

	if c.WebSocket != nil {
		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, d)
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
		}
		return
	}

	if "" != c.EventName {
		EventContent(c, d)
		return
	}

	_, writeError := c.Writer.Write(d)
	if writeError != nil {
		c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
	}
}

// Message sends utf-8 safe content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func Message(c *conn.Conn, m string) {
	Content(c, []byte(m))
}

// Messagef sends utf-8 safe content using a format.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func Messagef(c *conn.Conn, f string, v ...any) {
	Content(c, []byte(fmt.Sprintf(f, v...)))
}

// NotFound sends a message with status 404 Not Found.
func NotFound(c *conn.Conn, m string) {
	Status(c, http.StatusNotFound)
	Message(c, m)
}

// Unauthorized sends a message with status 401 Unauthorized.
func Unauthorized(c *conn.Conn, m string) {
	Status(c, http.StatusUnauthorized)
	Message(c, m)
}

// BadRequest sends a message with status 400 Bad Request.
func BadRequest(c *conn.Conn, m string) {
	Status(c, http.StatusBadRequest)
	Message(c, m)
}

// Error sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func Error(c *conn.Conn, e error) {
	Status(c, http.StatusBadRequest)
	Message(c, e.Error())
}

// Forbidden sends a message with status 403 Forbidden.
func Forbidden(c *conn.Conn, m string) {
	Status(c, http.StatusForbidden)
	Message(c, m)
}

// TooManyRequests sends a message with status 403 Forbidden.
func TooManyRequests(c *conn.Conn, m string) {
	Status(c, http.StatusTooManyRequests)
	Message(c, m)
}

// Json sends json content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func Json(c *conn.Conn, v any) {
	data, jsonError := json.Marshal(v)
	if jsonError != nil {
		c.Container.Config.ErrorLog.Println(jsonError, stack.Trace())
		return
	}

	if nil == c.WebSocket {
		contentType := c.Writer.Header().Get("Content-Type")
		if "" == contentType {
			c.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	Content(c, data)
}

// EmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func EmbeddedFileOrElse(c *conn.Conn, fs embed.FS, or func()) {
	fileName := c.Container.Config.PublicRoot + c.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !embeds.IsFile(fs, fileName) || embeds.IsDirectory(fs, fileName) {
		or()
		return
	}

	reader, readerInfo, readerError := embeds.NewFileReader(fs, fileName)
	if readerError != nil {
		c.Container.Config.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}
		return
	}

	if "" != c.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return
		}

		EventContent(c, data)
		return
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", mime.Parse(fileName))
	}

	if "" == c.Writer.Header().Get("Content-Length") {
		Header(c, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(c.Writer, c.Request, fileName, readerInfo.ModTime(), reader)
}

// FileOrElse sends the file requested by the client, or else falls back.
func FileOrElse(c *conn.Conn, or func()) {
	fileName := filepath.Join(c.Container.Config.PublicRoot, c.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		EmbeddedFileOrElse(c, c.Container.Config.Efs, or)
		return
	}

	reader, readerInfo, readerError := files.NewFileReader(fileName)
	if readerError != nil {
		c.Container.Config.ErrorLog.Println(readerError, stack.Trace())
		return
	}

	if c.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return
		}

		writeError := c.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	if "" != c.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return
		}

		EventContent(c, data)
		return
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", mime.Parse(fileName))
	}

	if "" == c.Writer.Header().Get("Content-Length") {
		Header(c, "Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(c.Writer, c.Request, fileName, readerInfo.ModTime(), reader)
}

// SseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SseUpgrade(c *conn.Conn) func(string) {
	Headers(c, map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Conn":                          "keep-alive",
	})

	c.EventName = "message"

	return func(eventName string) { c.EventName = eventName }
}

// WsUpgrade upgrades to web sockets.
func WsUpgrade(c *conn.Conn) {
	SseUpgradeWithUpgrader(c, websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// SseUpgradeWithUpgrader upgrades to web sockets.
func SseUpgradeWithUpgrader(c *conn.Conn, u websocket.Upgrader) {
	webSocketConnection, upgradeError := u.Upgrade(c.Writer, c.Request, nil)
	if upgradeError != nil {
		c.Container.Config.ErrorLog.Println(upgradeError, stack.Trace())
		return
	}

	defer func(webSocketConnection *websocket.Conn) {
		closeError := c.WebSocket.Close()
		if closeError != nil {
			c.Container.Config.ErrorLog.Println(closeError, stack.Trace())
		}
	}(webSocketConnection)

	c.WebSocket = webSocketConnection
	c.Locked = true

	return
}

// View sends a view.
func View(c *conn.Conn, v view.View) {
	if c.Writer.Header().Get("Location") != "" {
		return
	}

	if strings.Contains(c.Request.Header.Get("Accept"), "application/json") {
		if v.Data == nil {
			v.Data = map[string]any{}
		}
		props := map[string]any{
			"name":       v.Name,
			"data":       v.Data,
			"renderMode": v.RenderMode,
		}
		Json(c, props)
		return
	}

	html, err := view.Render(&v, c.Container)
	if err != nil {
		c.Container.Config.ErrorLog.Println(err, stack.Trace())
	}

	if "" == c.Writer.Header().Get("Content-Type") {
		Header(c, "Content-Type", "text/html")
	}

	Message(c, html)
}
