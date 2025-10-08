package send

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stacks"
)

// Header sends a header field.
//
// If the status has not been sent already, a default "200 OK" status will be sent immediately.
//
// This means the status will become locked and further attempts to send the status will fail with an error.
//
// All errors are sent to the server notifier.
func Header(client *clients.Client, key string, value string) {
	if client.Locked {
		client.Config.ErrorLog.Println("header is locked", stacks.Trace())
		return
	}

	client.Writer.Header().Set(key, value)
}

// Headers sends header fields.
func Headers(client *clients.Client, fields map[string]string) {
	if client.Locked {
		client.Config.ErrorLog.Println("header is locked", stacks.Trace())
		return
	}

	for key, value := range fields {
		client.Writer.Header().Set(key, value)
	}
}

// Redirect redirects the request to a location with a status.
func Redirect(client *clients.Client, location string, status int) {
	Status(client, status)
	Header(client, "Location", location)
}

// Navigate redirects the request to a location with status 302.
func Navigate(client *clients.Client, location string) {
	Redirect(client, location, 302)
	Message(client, "")
}

// Navigatef redirects the request to a location with status 302.
func Navigatef(client *clients.Client, format string, vars ...any) {
	Redirect(client, fmt.Sprintf(format, vars...), 302)
	Message(client, "")
}

// ContentType sets the Content-Type header field.
func ContentType(client *clients.Client, ctype string) {
	Header(client, "Content-Type", ctype)
}
