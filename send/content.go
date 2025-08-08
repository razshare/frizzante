package send

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"net/http"
)

// Flush send an empty message.
func Flush(c *client.Client) {
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
func Content(c *client.Client, d []byte) {
	if !c.Scope.Locked {
		c.Writer.WriteHeader(c.Scope.Status)
		c.Scope.Locked = true
	}

	if c.Scope.WebSocket != nil {
		writeError := c.Scope.WebSocket.WriteMessage(websocket.TextMessage, d)
		if writeError != nil {
			c.Scope.Container.Config.ErrorLog.Println(writeError, stack.Trace())
		}
		return
	}

	if "" != c.Scope.EventName {
		EventContent(c, d)
		return
	}

	_, writeError := c.Writer.Write(d)
	if writeError != nil {
		c.Scope.Container.Config.ErrorLog.Println(writeError, stack.Trace())
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
func Message(c *client.Client, m string) {
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
func Messagef(c *client.Client, f string, v ...any) {
	Content(c, []byte(fmt.Sprintf(f, v...)))
}

// NotFound sends a message with status 404 Not Found.
func NotFound(c *client.Client, m string) {
	Status(c, http.StatusNotFound)
	Message(c, m)
}

// Unauthorized sends a message with status 401 Unauthorized.
func Unauthorized(c *client.Client, m string) {
	Status(c, http.StatusUnauthorized)
	Message(c, m)
}

// BadRequest sends a message with status 400 Bad Request.
func BadRequest(c *client.Client, m string) {
	Status(c, http.StatusBadRequest)
	Message(c, m)
}

// Error sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func Error(c *client.Client, e error) {
	Status(c, http.StatusBadRequest)
	Message(c, e.Error())
}

// Forbidden sends a message with status 403 Forbidden.
func Forbidden(c *client.Client, m string) {
	Status(c, http.StatusForbidden)
	Message(c, m)
}

// TooManyRequests sends a message with status 403 Forbidden.
func TooManyRequests(c *client.Client, m string) {
	Status(c, http.StatusTooManyRequests)
	Message(c, m)
}
