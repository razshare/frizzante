package send

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Flush send an empty message.
//
// Compatible with web sockets and server sent events.
func Flush(http *scopes.Http) {
	Message(http, "")
}
