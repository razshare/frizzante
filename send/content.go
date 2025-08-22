package send

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"net/http"
)

// Flush send an empty message.
func Flush(client *client.Client) {
	Message(client, "")
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
func Content(client *client.Client, d []byte) {
	if !client.Locked {
		client.Writer.WriteHeader(client.Status)
		client.Locked = true
	}

	if client.WebSocket != nil {
		if err := client.WebSocket.WriteMessage(websocket.TextMessage, d); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
		}
		return
	}

	if "" != client.EventName {
		EventContent(client, d)
		return
	}

	if _, err := client.Writer.Write(d); err != nil {
		client.Config.ErrorLog.Println(err, stack.Trace())
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
func Message(client *client.Client, m string) {
	Content(client, []byte(m))
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
func Messagef(client *client.Client, f string, v ...any) {
	Content(client, []byte(fmt.Sprintf(f, v...)))
}

// NotFound sends a message with status 404 Not Found.
func NotFound(client *client.Client, m string) {
	Status(client, http.StatusNotFound)
	Message(client, m)
}

// Unauthorized sends a message with status 401 Unauthorized.
func Unauthorized(client *client.Client, m string) {
	Status(client, http.StatusUnauthorized)
	Message(client, m)
}

// BadRequest sends a message with status 400 Bad Request.
func BadRequest(client *client.Client, m string) {
	Status(client, http.StatusBadRequest)
	Message(client, m)
}

// Error sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func Error(client *client.Client, e error) {
	Status(client, http.StatusBadRequest)
	Message(client, e.Error())
}

// Forbidden sends a message with status 403 Forbidden.
func Forbidden(client *client.Client, m string) {
	Status(client, http.StatusForbidden)
	Message(client, m)
}

// TooManyRequests sends a message with status 403 Forbidden.
func TooManyRequests(client *client.Client, m string) {
	Status(client, http.StatusTooManyRequests)
	Message(client, m)
}
