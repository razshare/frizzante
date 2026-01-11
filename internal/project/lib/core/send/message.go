package send

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Message sends utf-8 safe content.
//
// Compatible with web sockets and server sent events.
func Message(client *clients.Client, message string) {
	Content(client, []byte(message))
}
