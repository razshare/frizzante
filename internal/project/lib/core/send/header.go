package send

import (
	"main/lib/core/client"

	"github.com/razshare/frizzante/stack"
)

// Header sends a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func Header(client *client.Client, key string, value string) {
	if client.Locked {
		client.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	client.Writer.Header().Set(key, value)
}

// Headers sends header fields.
func Headers(client *client.Client, fields map[string]string) {
	if client.Locked {
		client.Config.ErrorLog.Println("header is locked", stack.Trace())
		return
	}

	for key, value := range fields {
		client.Writer.Header().Set(key, value)
	}
}

// Redirect redirects the request to a location with a status.
func Redirect(client *client.Client, location string, status int) {
	Status(client, status)
	Header(client, "Location", location)
}

// Navigate redirects the request to a location with status 302.
func Navigate(client *client.Client, location string) {
	Redirect(client, location, 302)
	Message(client, "")
}

// ContentType sets the Content-Type header field.
func ContentType(client *client.Client, ctype string) {
	Header(client, "Content-Type", ctype)
}
