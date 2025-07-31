package servers

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/mimes"
	"github.com/razshare/frizzante/traces"
	"github.com/razshare/frizzante/views"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

func New() *Server {
	return &Server{
		Connections:   map[string]*net.Conn{},
		InfoLog:       log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
		Address:       "0.0.0.0:8080",
		SecureAddress: "0.0.0.0:8383",
		PublicRoot:    "app/dist/client",
		AppRoot:       "app",
		ServerJs:      "app/dist/server.js",
		IndexHtml:     "app/dist/client/index.html",
		Server: http.Server{
			Handler:        http.NewServeMux(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 3 * globals.MB,
			ErrorLog:       log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		},
	}
}

// Start starts the server.
//
// If the server fails to start, ServerStart crashes the program.
func (web *Server) Start() {
	mux := web.Handler.(*http.ServeMux)

	for _, route := range web.Routes {
		mux.HandleFunc(route.Pattern, func(writer http.ResponseWriter, request *http.Request) {
			con := &Connection{
				Web:     web,
				Request: request,
				Writer:  writer,
				Status:  200,
				EventId: 1,
			}

			for _, tag := range route.Tags {
				for _, guard := range web.Guards {
					if !slices.Contains(guard.Tags, tag) {
						continue
					}
					allowed := false
					guard.Handler(con, func() { allowed = true })
					if !allowed {
						web.InfoLog.Printf("route `%s` tagged with `%s` denied the request because guard `%s` did not pass", route.Pattern, tag, guard.Name)
						return
					}
				}
			}

			route.Handler(con)
		})
	}

	var group sync.WaitGroup

	group.Add(2)

	go func() {
		web.InfoLog.Printf("listening for requests at http://%s", web.Address)
		serveError := http.ListenAndServe(web.Address, web.Handler)
		if serveError != nil {
			if errors.Is(serveError, http.ErrServerClosed) {
				web.InfoLog.Println("shutting down server")
				return
			}
			log.Fatal(serveError)
		}
	}()

	go func() {
		if "" != web.Certificate && "" != web.Key {
			web.InfoLog.Printf("listening for requests at https://%s", web.SecureAddress)
			serveError := http.ListenAndServeTLS(web.SecureAddress, web.Certificate, web.Key, web.Handler)
			if serveError != nil {
				if errors.Is(serveError, http.ErrServerClosed) {
					web.InfoLog.Printf("shutting down server")
					return
				}
				log.Fatal(serveError)
			}
		}
	}()

	group.Wait()
}

// Stop attempts to stop the server.
//
// If the shutdown attempt fails, ServerStop crashes the program.
func (web *Server) Stop() {
	if err := web.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
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
	return connection.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func (connection *Connection) IsAlive() *bool {
	isAlive := true
	go func() {
		<-connection.ReceiveCancellation()
		isAlive = false
	}()
	return &isAlive
}

// ReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveCookie(key string) string {
	cookie, cookieError := connection.Request.Cookie(key)
	if cookieError != nil {
		traces.Trace(connection.Web.ErrorLog, cookieError)
		return ""
	}

	data, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		traces.Trace(connection.Web.ErrorLog, queryError)
		return ""
	}

	return data
}

// ReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveMessage() string {
	if connection.WebSocket != nil {
		_, data, readError := connection.WebSocket.ReadMessage()
		if readError != nil {
			traces.Trace(connection.Web.ErrorLog, readError)
			return ""
		}
		return string(data)
	}

	data, readError := io.ReadAll(connection.Request.Body)
	if readError != nil {
		traces.Trace(connection.Web.ErrorLog, readError)
		return ""
	}
	return string(data)
}

// ReceiveJson reads the next JSON-encoded message from the
// connection and stores it in the value pointed to by val.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveJson(value any) {
	if connection.WebSocket != nil {
		jsonError := connection.WebSocket.ReadJSON(value)
		if jsonError != nil {
			traces.Trace(connection.Web.ErrorLog, jsonError)
			return
		}
		return
	}

	data, readError := io.ReadAll(connection.Request.Body)
	if readError != nil {
		traces.Trace(connection.Web.ErrorLog, readError)
		return
	}

	jsonError := json.Unmarshal(data, value)
	if jsonError != nil {
		traces.Trace(connection.Web.ErrorLog, jsonError)
		return
	}
}

// ReceiveForm reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of 2MB
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func (connection *Connection) ReceiveForm() url.Values {
	return connection.ReceiveFormWithMaxMemory(2 * globals.MB)
}

// ReceiveFormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of maxMemory bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func (connection *Connection) ReceiveFormWithMaxMemory(maxMemory int64) url.Values {
	if connection.WebSocket != nil {
		traces.Trace(connection.Web.ErrorLog, errors.New("connection is not of type web socket"))
		return url.Values{}
	}

	formError := connection.Request.ParseMultipartForm(maxMemory)
	if formError != nil {
		if !errors.Is(formError, http.ErrNotMultipart) {
			return url.Values{}
		}

		formError = connection.Request.ParseForm()
		if formError != nil {
			traces.Trace(connection.Web.ErrorLog, formError)
			return url.Values{}
		}
	}

	return connection.Request.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceiveQuery(key string) string {
	return connection.Request.URL.Query().Get(key)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func (connection *Connection) ReceivePath(key string) string {
	return connection.Request.PathValue(key)
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
func (connection *Connection) ReceiveContentType() string {
	return connection.Request.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func (connection *Connection) VerifyContentType(contentType ...string) bool {
	requestedMime := connection.Request.Header.Get("Content-Type")
	for _, acceptedMime := range contentType {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func (connection *Connection) VerifyAccept(accept ...string) bool {
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
func (connection *Connection) SendEventContent(content []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", connection.EventId, connection.EventName)

	_, writeError := connection.Writer.Write([]byte(header))
	if writeError != nil {
		traces.Trace(connection.Web.ErrorLog, writeError)
		return
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeError = connection.Writer.Write([]byte("data: "))
		if writeError != nil {
			traces.Trace(connection.Web.ErrorLog, writeError)
			return
		}

		_, writeError = connection.Writer.Write(line)
		if writeError != nil {
			traces.Trace(connection.Web.ErrorLog, writeError)
			return
		}

		_, writeError = connection.Writer.Write([]byte("\r\n"))
		if writeError != nil {
			traces.Trace(connection.Web.ErrorLog, writeError)
			return
		}
	}

	_, writeError = connection.Writer.Write([]byte("\r\n"))
	if writeError != nil {
		traces.Trace(connection.Web.ErrorLog, writeError)
		return
	}

	flusher, flushedOk := connection.Writer.(http.Flusher)
	if !flushedOk {
		traces.Trace(connection.Web.ErrorLog, errors.New("could not retrieve flusher"))
		return
	}

	flusher.Flush()

	connection.EventId++
}

// SendNavigate redirects the request to a location with status 302.
func (connection *Connection) SendNavigate(location string) {
	connection.SendRedirect(location, 302)
	connection.SendFlush()
}

// SendRedirect redirects the request to a location with a status.
func (connection *Connection) SendRedirect(location string, status int) {
	connection.SendStatus(status)
	connection.SendHeader("Location", location)
}

// SendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func (connection *Connection) SendStatus(status int) {
	if connection.Locked {
		traces.Trace(connection.Web.ErrorLog, "status is locked")
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
func (connection *Connection) SendHeader(key string, value string) {
	if connection.Locked {
		traces.Trace(connection.Web.ErrorLog, "header is locked")
		return
	}

	connection.Writer.Header().Set(key, value)
}

func (connection *Connection) SendHeaders(headers map[string]string) {
	if connection.Locked {
		traces.Trace(connection.Web.ErrorLog, "header is locked")
		return
	}

	for key, value := range headers {
		connection.Writer.Header().Set(key, value)
	}
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
		if writeError != nil {
			traces.Trace(connection.Web.ErrorLog, writeError)
		}
		return
	}

	if "" != connection.EventName {
		connection.SendEventContent(content)
		return
	}

	_, writeError := connection.Writer.Write(content)
	if writeError != nil {
		traces.Trace(connection.Web.ErrorLog, writeError)
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
func (connection *Connection) SendMessage(value string) {
	connection.SendContent([]byte(value))
}

// SendNotFound sends a message with status 404 Not Found.
func (connection *Connection) SendNotFound(value string) {
	connection.SendStatus(http.StatusNotFound)
	connection.SendMessage(value)
}

// SendUnauthorized sends a message with status 401 Unauthorized.
func (connection *Connection) SendUnauthorized(value string) {
	connection.SendStatus(http.StatusUnauthorized)
	connection.SendMessage(value)
}

// SendBadRequest sends a message with status 400 Bad Request.
func (connection *Connection) SendBadRequest(value string) {
	connection.SendStatus(http.StatusBadRequest)
	connection.SendMessage(value)
}

// SendError sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func (connection *Connection) SendError(err error) {
	connection.SendStatus(http.StatusBadRequest)
	connection.SendMessage(err.Error())
}

// SendForbidden sends a message with status 403 Forbidden.
func (connection *Connection) SendForbidden(v string) {
	connection.SendStatus(http.StatusForbidden)
	connection.SendMessage(v)
}

// SendTooManyRequests sends a message with status 403 Forbidden.
func (connection *Connection) SendTooManyRequests(v string) {
	connection.SendStatus(http.StatusTooManyRequests)
	connection.SendMessage(v)
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
func (connection *Connection) SendJson(value any) {
	data, jsonError := json.Marshal(value)
	if jsonError != nil {
		traces.Trace(connection.Web.ErrorLog, jsonError)
		return
	}

	if nil == connection.WebSocket {
		contentType := connection.Writer.Header().Get("Content-Type")
		if "" == contentType {
			connection.Writer.Header().Set("Content-Type", "application/json")
		}
	}

	connection.SendContent(data)
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func (connection *Connection) SendEmbeddedFileOrElse(efs embed.FS, orElse func()) {
	fileName := connection.Web.PublicRoot + connection.Request.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !embeds.IsFile(efs, fileName) || embeds.IsDirectory(efs, fileName) {
		orElse()
		return
	}

	reader, readerInfo, readerError := embeds.FileReader(efs, fileName)
	if readerError != nil {
		traces.Trace(connection.Web.ErrorLog, readerError)
		return
	}

	if connection.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(connection.Web.ErrorLog, readError)
			return
		}

		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			traces.Trace(connection.Web.ErrorLog, writeError)
			return
		}
		return
	}

	if "" != connection.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(connection.Web.ErrorLog, readError)
			return
		}

		connection.SendEventContent(data)
		return
	}

	if "" == connection.Writer.Header().Get("Content-Type") {
		connection.SendHeader("Content-Type", mimes.Mime(fileName))
	}

	if "" == connection.Writer.Header().Get("Content-Length") {
		connection.SendHeader("Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(connection.Writer, connection.Request, fileName, readerInfo.ModTime(), reader)
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func (connection *Connection) SendFileOrElse(orElse func()) {
	fileName := filepath.Join(connection.Web.PublicRoot, connection.Request.RequestURI)

	if !files.IsFile(fileName) || files.IsDirectory(fileName) {
		connection.SendEmbeddedFileOrElse(connection.Web.Efs, orElse)
		return
	}

	reader, readerInfo, readerError := files.FileReader(fileName)
	if readerError != nil {
		traces.Trace(connection.Web.ErrorLog, readerError)
		return
	}

	if connection.WebSocket != nil {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(connection.Web.ErrorLog, readError)
			return
		}

		writeError := connection.WebSocket.WriteMessage(websocket.TextMessage, data)
		if writeError != nil {
			traces.Trace(connection.Web.ErrorLog, writeError)
			return
		}
	}

	if "" != connection.EventName {
		data, readError := io.ReadAll(reader)
		if readError != nil {
			traces.Trace(connection.Web.ErrorLog, readError)
			return
		}

		connection.SendEventContent(data)
		return
	}

	if "" == connection.Writer.Header().Get("Content-Type") {
		connection.SendHeader("Content-Type", mimes.Mime(fileName))
	}

	if "" == connection.Writer.Header().Get("Content-Length") {
		connection.SendHeader("Content-Length", fmt.Sprintf("%d", readerInfo.Size()))
	}

	http.ServeContent(connection.Writer, connection.Request, fileName, readerInfo.ModTime(), reader)
}

// SendSseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func (connection *Connection) SendSseUpgrade() func(eventName string) {
	connection.SendHeaders(map[string]string{
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
func (connection *Connection) SendWsUpgrade() {
	connection.SendConfiguredWsUpgrade(websocket.Upgrader{
		ReadBufferSize:  10 * globals.KB,
		WriteBufferSize: 10 * globals.KB,
	})
}

// SendConfiguredWsUpgrade upgrades to web sockets.
func (connection *Connection) SendConfiguredWsUpgrade(upgrader websocket.Upgrader) {
	webSocketConnection, upgradeError := upgrader.Upgrade(connection.Writer, connection.Request, nil)
	if upgradeError != nil {
		traces.Trace(connection.Web.ErrorLog, upgradeError)
		return
	}

	defer func(webSocketConnection *websocket.Conn) {
		closeError := webSocketConnection.Close()
		if closeError != nil {
			traces.Trace(connection.Web.ErrorLog, closeError)
		}
	}(webSocketConnection)

	connection.WebSocket = webSocketConnection
	connection.Locked = true

	return
}

// SendView sends a view.
func (connection *Connection) SendView(view views.View) {
	if connection.Writer.Header().Get("Location") != "" {
		return
	}

	if connection.VerifyAccept("application/json") {
		if view.Data == nil {
			view.Data = map[string]any{}
		}
		props := map[string]any{
			"name":       view.Name,
			"data":       view.Data,
			"renderMode": view.RenderMode,
		}
		connection.SendJson(props)
		return
	}

	if view.ServerJs == "" {
		view.ServerJs = connection.Web.ServerJs
	}

	if view.IndexHtml == "" {
		view.IndexHtml = connection.Web.IndexHtml
	}

	if view.AppRoot == "" {
		view.AppRoot = connection.Web.AppRoot
	}

	html, renderError := view.Render(connection.Web.Efs)
	if renderError != nil {
		traces.Trace(connection.Web.ErrorLog, renderError)
	}

	if "" == connection.Writer.Header().Get("Content-Type") {
		connection.SendHeader("Content-Type", "text/html")
	}

	connection.SendMessage(html)
}
