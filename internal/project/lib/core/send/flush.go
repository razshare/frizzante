package send

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Flush send an empty message.
//
// Compatible with web sockets and server sent events.
func Flush(client *clients.Client) {
	Message(client, "")
}
