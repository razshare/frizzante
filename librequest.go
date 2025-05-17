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
	server        *Server
	response      *Response
	httpRequest   *http.Request
	webSocketConn *websocket.Conn
}

// RequestReceiveCancellation returns a channel that's closed when the request is cancelled.
func RequestReceiveCancellation(self *Request) <-chan struct{} {
	return self.httpRequest.Context().Done()
}

// RequestReceiveCookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func RequestReceiveCookie(self *Request, key string) string {
	cookie, cookieError := self.httpRequest.Cookie(key)
	if nil != cookieError {
		NotifierSendError(self.server.notifier, cookieError)
		return ""
	}
	value, unescapeError := url.QueryUnescape(cookie.Value)
	if nil != unescapeError {
		return ""
	}

	return value
}

// RequestReceiveMessage reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func RequestReceiveMessage(self *Request) string {
	if self.webSocketConn != nil {
		_, readBytes, readError := self.webSocketConn.ReadMessage()
		if nil != readError {
			NotifierSendError(self.server.notifier, readError)
			return ""
		}
		return string(readBytes)
	}

	readBytes, readAllError := io.ReadAll(self.httpRequest.Body)
	if nil != readAllError {
		NotifierSendError(self.server.notifier, readAllError)
		return ""
	}
	return string(readBytes)
}

// RequestReceiveJson reads the message as json and returns the value and a boolean,
// which indicates success or failure.
//
// Compatible with web sockets.
func RequestReceiveJson[T any](self *Request) *T {
	var value T
	if self.webSocketConn != nil {
		jsonError := self.webSocketConn.ReadJSON(value)
		if nil != jsonError {
			NotifierSendError(self.server.notifier, jsonError)
			return nil
		}
		return &value
	}

	readBytes, readAllError := io.ReadAll(self.httpRequest.Body)
	if nil != readAllError {
		NotifierSendError(self.server.notifier, readAllError)
		return nil
	}
	unmarshalError := json.Unmarshal(readBytes, &value)
	if nil != unmarshalError {
		NotifierSendError(self.server.notifier, unmarshalError)
		return nil
	}
	return &value
}

// RequestReceiveForm reads the message as a form and returns the value.
func RequestReceiveForm(self *Request) *url.Values {
	if self.webSocketConn != nil {
		NotifierSendError(self.server.notifier, errors.New("web socket connections cannot receive form payloads"))
		return &url.Values{}
	}

	parseMultipartFormError := self.httpRequest.ParseMultipartForm(self.server.multipartFormMaxMemory)
	if nil != parseMultipartFormError {
		if !errors.Is(parseMultipartFormError, http.ErrNotMultipart) {
			NotifierSendError(self.server.notifier, parseMultipartFormError)
		}

		parseFormError := self.httpRequest.ParseForm()
		if nil != parseFormError {
			NotifierSendError(self.server.notifier, parseFormError)
		}
	}

	return &self.httpRequest.Form
}

// RequestReceiveQuery reads a query field and returns the value.
//
// Compatible with web sockets.
func RequestReceiveQuery(self *Request, name string) string {
	return self.httpRequest.URL.Query().Get(name)
}

// RequestReceivePath reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func RequestReceivePath(self *Request, name string) string {
	return self.httpRequest.PathValue(name)
}

// RequestReceiveHeader reads a header field and returns the value.
//
// Compatible with web sockets.
func RequestReceiveHeader(self *Request, key string) string {
	return self.httpRequest.Header.Get(key)
}

// RequestReceiveContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func RequestReceiveContentType(self *Request) string {
	return self.httpRequest.Header.Get("Content-Type")
}

// RequestVerifyContentType checks if the incoming request has any of the given content-types.
func RequestVerifyContentType(self *Request, contentTypes ...string) bool {
	requestedMime := self.httpRequest.Header.Get("Content-Type")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.HasPrefix(requestedMime, acceptedMime) {
			return true
		}
	}

	return false
}

// RequestVerifyAccept checks if the incoming request accepts any of the given content-types.
func RequestVerifyAccept(self *Request, contentTypes ...string) bool {
	requestedAcceptMime := self.httpRequest.Header.Get("Accept")
	for _, acceptedMime := range contentTypes {
		if acceptedMime == "*" || strings.Contains(requestedAcceptMime, acceptedMime) {
			return true
		}
	}

	return false
}
