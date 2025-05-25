package frizzante

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

type Response struct {
	server     *Server
	request    *Request
	writer     *http.ResponseWriter
	locked     bool
	statusCode int
	header     *http.Header
	webSocket  *websocket.Conn
	eventName  string
	eventId    int64
	context    map[string]any
}

type Navigate struct {
	Location string
}

// responseSendEventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a Server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func (response *Response) sendEventContent(content []byte) *Response {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", response.eventId, response.eventName)

	_, writeEventError := (*response.writer).Write([]byte(header))
	if nil != writeEventError {
		response.server.notifier.SendError(writeEventError)
		return response
	}

	for _, line := range bytes.Split(content, []byte("\r\n")) {
		_, writeEventError = (*response.writer).Write([]byte("data: "))
		if nil != writeEventError {
			response.server.notifier.SendError(writeEventError)
			return response
		}

		_, writeEventError = (*response.writer).Write(line)
		if nil != writeEventError {
			response.server.notifier.SendError(writeEventError)
			return response
		}

		_, writeEventError = (*response.writer).Write([]byte("\r\n"))
		if nil != writeEventError {
			response.server.notifier.SendError(writeEventError)
			return response
		}
	}

	_, writeEventError = (*response.writer).Write([]byte("\r\n"))
	if nil != writeEventError {
		response.server.notifier.SendError(writeEventError)
		return response
	}

	flusher, flushedOk := (*response.writer).(http.Flusher)
	if !flushedOk {
		response.server.notifier.SendError(errors.New("could not retrieve flusher"))
		return response
	}

	flusher.Flush()

	response.eventId++
	return response
}

// SendNavigate redirects the request with status 302.
func (response *Response) SendNavigate(location string) *Response {
	response.SendRedirect(location, 302)
	response.SendMessage("")
	return response
}

// SendRedirect redirects the request.
func (response *Response) SendRedirect(location string, statusCode int) *Response {
	response.SendStatus(statusCode)
	response.SendHeader("Location", location)
	return response
}

// SendStatus sets the status code.
//
// This will lock the status, which makes it
// so that the increaseIndex time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func (response *Response) SendStatus(code int) *Response {
	if response.locked {
		response.server.notifier.SendError(errors.New("status is locked"))
		return response
	}
	response.statusCode = code
	return response
}

// SendHeader sets a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func (response *Response) SendHeader(key string, value string) *Response {
	if response.locked {
		response.server.notifier.SendError(errors.New("headers locked"))
		return response
	}

	response.header.Set(key, value)
	return response
}

// SendContentType sets the Content-Type header field.
func (response *Response) SendContentType(contentType string) *Response {
	response.SendHeader("Content-Type", contentType)
	return response
}

// SendCookie sends a cookies to the client.
func (response *Response) SendCookie(key string, value string) *Response {
	response.SendHeader("Set-Cookie", fmt.Sprintf("%s=%s; Path=/; HttpOnly", url.QueryEscape(key), url.QueryEscape(value)))
	return response
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
func (response *Response) SendContent(content []byte) *Response {
	if !response.locked {
		(*response.writer).WriteHeader(response.statusCode)
		response.locked = true
	}

	if response.webSocket != nil {
		writeError := response.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			response.server.notifier.SendError(writeError)
			return response
		}
		return response
	}

	if "" != response.eventName {
		response.sendEventContent(content)
		return response
	}

	_, err := (*response.writer).Write(content)
	if nil != err {
		response.server.notifier.SendError(err)
		return response
	}
	return response
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
func (response *Response) SendMessage(message string) *Response {
	response.SendContent([]byte(message))
	return response
}

// SendNotFound sends a message with status 404 Not Found.
func (response *Response) SendNotFound(message string) *Response {
	response.SendStatus(http.StatusNotFound)
	response.SendMessage(message)
	return response
}

// SendUnauthorized sends a message with status 401 Unauthorized.
func (response *Response) SendUnauthorized(message string) *Response {
	response.SendStatus(http.StatusUnauthorized)
	response.SendMessage(message)
	return response
}

// SendBadRequest sends a message with status 400 Bad Request.
func (response *Response) SendBadRequest(message string) *Response {
	response.SendStatus(http.StatusBadRequest)
	response.SendMessage(message)
	return response
}

// SendInternalServerError sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func (response *Response) SendInternalServerError(err error) *Response {
	response.server.notifier.SendError(err)
	response.SendStatus(http.StatusBadRequest)
	response.SendMessage(err.Error())
	return response
}

// SendForbidden sends a message with status 403 Forbidden.
func (response *Response) SendForbidden(message string) *Response {
	response.SendStatus(http.StatusForbidden)
	response.SendMessage(message)
	return response
}

// SendTooManyRequests sends a message with status 403 Forbidden.
func (response *Response) SendTooManyRequests(message string) *Response {
	response.SendStatus(http.StatusTooManyRequests)
	response.SendMessage(message)
	return response
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
func (response *Response) SendJson(payload any) *Response {
	content, marshalError := json.Marshal(payload)
	if nil != marshalError {
		response.server.notifier.SendError(marshalError)
		return response
	}

	if nil == response.webSocket {
		contentType := response.header.Get("Content-Type")
		if "" == contentType {
			response.header.Set("Content-Type", "application/json")
		}
	}

	response.SendContent(content)
	return response
}

// SendEmbeddedFileOrElse sends the embedded file requested by the client,
// or the closest index.html embedded file, or else falls back.
func (response *Response) SendEmbeddedFileOrElse(orElse func()) *Response {
	request := response.request
	fileName := ".dist/client" + request.httpRequest.RequestURI
	fileName = strings.Split(fileName, "?")[0]
	fileName = strings.Split(fileName, "&")[0]

	if !existsInEmbeddedFileSystem(request.server.dist, fileName) ||
		isEmbeddedDirectory(request.server.dist, fileName) {
		orElse()
		return response
	}

	reader, info, readerError := createReaderFromEmbeddedFileName(
		request.server.dist,
		fileName,
	)
	if nil != readerError {
		response.server.notifier.SendError(readerError)
		return response
	}

	if response.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			response.server.notifier.SendError(readError)
			return response
		}
		writeError := response.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			response.server.notifier.SendError(writeError)
		}
		return response
	}

	if "" != response.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			response.server.notifier.SendError(readError)
			return response
		}
		response.sendEventContent(content)
		return response
	}

	if "" == response.header.Get("Content-Type") {
		response.SendHeader("Content-Type", mime(fileName))
	}

	if "" == response.header.Get("Content-Length") {
		response.SendHeader("Content-Length", fmt.Sprintf("%d", (*info).Size()))
	}
	http.ServeContent(*response.writer, request.httpRequest, fileName, (*info).ModTime(), reader)
	return response
}

// SendFileOrElse sends the file requested by the client, or else falls back.
func (response *Response) SendFileOrElse(orElse func()) *Response {
	request := response.request
	fileName := filepath.Join(".dist", "client", request.httpRequest.RequestURI)

	if !fileExists(fileName) || isDirectory(fileName) {
		response.SendEmbeddedFileOrElse(orElse)
		return response
	}

	reader, info, readerError := createReaderFromFileName(fileName)
	if nil != readerError {
		response.server.notifier.SendError(readerError)
		return response
	}

	if response.webSocket != nil {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			response.server.notifier.SendError(readError)
			return response
		}
		writeError := response.webSocket.WriteMessage(websocket.TextMessage, content)
		if nil != writeError {
			response.server.notifier.SendError(writeError)
		}
		return response
	}

	if "" != response.eventName {
		content, readError := io.ReadAll(reader)
		if nil != readError {
			response.server.notifier.SendError(readError)
			return response
		}
		response.sendEventContent(content)
		return response
	}

	if "" == response.header.Get("Content-Type") {
		response.SendHeader("Content-Type", mime(fileName))
	}

	if "" == response.header.Get("Content-Length") {
		response.SendHeader("Content-Length", fmt.Sprintf("%d", (*info).Size()))
	}
	http.ServeContent(*response.writer, request.httpRequest, fileName, (*info).ModTime(), reader)
	return response
}

// SendSseUpgrade upgrades the http connection to server sent events
// and returns a function that sets the name of the current event.
//
// The default event is "message".
func (response *Response) SendSseUpgrade() (setEventName func(eventName string)) {
	response.SendHeader("Access-Control-Allow-Origin", "*")
	response.SendHeader("Access-Control-Expose-Headers", "Content-Type")
	response.SendHeader("Content-Type", "text/event-stream")
	response.SendHeader("Cache-Control", "no-cache")
	response.SendHeader("Connection", "keep-alive")
	response.eventName = "message"
	setEventName = func(eventName string) {
		if "" == eventName {
			response.server.notifier.SendError(
				fmt.Errorf("renaming a server sent event (`%s`) to an empty string is not allowed", response.eventName),
			)
			return
		}

		response.eventName = eventName
	}
	return
}

// SendWsUpgrade upgrades the http connection to web sockets.
func (response *Response) SendWsUpgrade() *Response {
	request := response.request
	conn, upgradeError := response.server.upgrader.Upgrade(*response.writer, request.httpRequest, nil)
	if nil != upgradeError {
		request.server.notifier.SendError(upgradeError)
		return response
	}
	defer func(conn *websocket.Conn) {
		closeError := conn.Close()
		if nil != closeError {
			request.server.notifier.SendError(closeError)
		}
	}(conn)
	response.webSocket = conn
	request.webSocket = conn
	response.locked = true
	return response
}

// SendView sends a view.
func (response *Response) SendView(view View) *Response {
	if "" != response.header.Get("Location") {
		return response
	}

	if nil == view.Data {
		view.Data = map[string]string{}
	}

	if response.request.VerifyAccept("application/json") {
		response.SendJson(&ServerProperties{
			View:       view.Name,
			RenderMode: view.RenderMode,
			Data:       view.Data,
		})
		return response
	}

	content, compileError := view.Render(response.server.dist)
	if nil != compileError {
		response.server.notifier.SendError(compileError)
		return response
	}

	if "" == response.header.Get("Content-Type") {
		response.SendHeader("Content-Type", "text/html")
	}

	response.SendMessage(content)
	return response
}
