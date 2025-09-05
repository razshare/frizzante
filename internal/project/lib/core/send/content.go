package send

import (
	"fmt"
	"main/lib/core/client"
	"net/http"

	"github.com/gorilla/websocket"
	"main/lib/core/stack"
)

// Content sends binary safe content.
//
// If the status code or the header have not been sent already, a default status of "200 OK" will be sent immediately along with whatever headers you've previously defined.
//
// The status code and the header will become locked and further attempts to send either of them will fail with an error.
//
// All errors are sent to the server notifier.
//
// Compatible with web sockets.
func Content(client *client.Client, data []byte) {
	if !client.Locked {
		client.Writer.WriteHeader(client.Status)
		client.Locked = true
	}

	if client.WebSocket != nil {
		if err := client.WebSocket.WriteMessage(websocket.TextMessage, data); err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
		}
		return
	}

	if "" != client.EventName {
		EventContent(client, data)
		return
	}

	if _, err := client.Writer.Write(data); err != nil {
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
func Message(client *client.Client, message string) {
	Content(client, []byte(message))
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
func Messagef(client *client.Client, format string, vars ...any) {
	Content(client, []byte(fmt.Sprintf(format, vars...)))
}

// NotFound sends a message with status 404 Not Found.
func NotFound(client *client.Client, message string) {
	Status(client, http.StatusNotFound)
	Message(client, message)
}

// Unauthorized sends a message with status 401 Unauthorized.
func Unauthorized(client *client.Client, message string) {
	Status(client, http.StatusUnauthorized)
	Message(client, message)
}

// BadRequest sends a message with status 400 Bad Request.
func BadRequest(client *client.Client, message string) {
	Status(client, http.StatusBadRequest)
	Message(client, message)
}

// Error sends a message with status 500 Internal server Error
// and also sends the error to the server notifier.
func Error(client *client.Client, err error) {
	Status(client, http.StatusInternalServerError)
	Message(client, err.Error())
}

// Forbidden sends a message with status 403 Forbidden.
func Forbidden(client *client.Client, message string) {
	Status(client, http.StatusForbidden)
	Message(client, message)
}

// TooManyRequests sends a message with status 403 Forbidden.
func TooManyRequests(client *client.Client, message string) {
	Status(client, http.StatusTooManyRequests)
	Message(client, message)
}

// Flush send an empty message.
func Flush(client *client.Client) {
	Message(client, "")
}
