package frizzante

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	server      *Server
	response    *Response
	httpRequest *http.Request
	webSocket   *websocket.Conn
}

// ReceiveCancellation returns a channel that closes when the request gets cancelled.
func (request *Request) ReceiveCancellation() <-chan struct{} {
	return request.httpRequest.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func (request *Request) IsAlive() *bool {
	value := true
	go func() {
		<-request.ReceiveCancellation()
		value = false
	}()
	return &value
}

// ReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func (request *Request) ReceiveCookie(key string) string {
	cookie, cookieError := request.httpRequest.Cookie(key)
	if nil != cookieError {
		request.server.notifier.SendError(cookieError)
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
func (request *Request) ReceiveMessage() string {
	if request.webSocket != nil {
		_, readBytes, readError := request.webSocket.ReadMessage()
		if nil != readError {
			request.server.notifier.SendError(readError)
			return ""
		}
		return string(readBytes)
	}

	readBytes, readAllError := io.ReadAll(request.httpRequest.Body)
	if nil != readAllError {
		request.server.notifier.SendError(readAllError)
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
func (request *Request) ReceiveJson(out any) bool {
	if request.webSocket != nil {
		jsonError := request.webSocket.ReadJSON(out)
		if nil != jsonError {
			request.server.notifier.SendError(jsonError)
			return false
		}
		return true
	}

	readBytes, readAllError := io.ReadAll(request.httpRequest.Body)
	if nil != readAllError {
		request.server.notifier.SendError(readAllError)
		return false
	}
	unmarshalError := json.Unmarshal(readBytes, out)
	if nil != unmarshalError {
		request.server.notifier.SendError(unmarshalError)
		return false
	}
	return true
}

// ReceiveForm reads the message as a form and returns the value.
func (request *Request) ReceiveForm() *url.Values {
	if request.webSocket != nil {
		request.server.notifier.SendError(errors.New("web socket connections cannot receive form payloads"))
		return &url.Values{}
	}

	parseMultipartFormError := request.httpRequest.ParseMultipartForm(request.server.formMaxMemory)
	if nil != parseMultipartFormError {
		if !errors.Is(parseMultipartFormError, http.ErrNotMultipart) {
			request.server.notifier.SendError(parseMultipartFormError)
		}

		parseFormError := request.httpRequest.ParseForm()
		if nil != parseFormError {
			request.server.notifier.SendError(parseFormError)
		}
	}

	return &request.httpRequest.Form
}

// ReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func (request *Request) ReceiveQuery(name string) string {
	return request.httpRequest.URL.Query().Get(name)
}

// ReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func (request *Request) ReceivePath(name string) string {
	return request.httpRequest.PathValue(name)
}

// ReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func (request *Request) ReceiveHeader(key string) string {
	return request.httpRequest.Header.Get(key)
}

// ReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func (request *Request) ReceiveContentType() string {
	return request.httpRequest.Header.Get("Content-Type")
}

// VerifyContentType checks if the incoming request has any of the given content-types.
func (request *Request) VerifyContentType(contentTypes ...string) bool {
	requestedMime := request.httpRequest.Header.Get("Content-Type")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// VerifyAccept checks if the incoming request accepts any of the given content-types.
func (request *Request) VerifyAccept(contentTypes ...string) bool {
	requestedAcceptMime := request.httpRequest.Header.Get("Accept")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.Contains(requestedAcceptMime, acceptedMime) {
			return true
		}
	}

	return false
}
