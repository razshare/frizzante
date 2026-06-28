package send

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Message sends utf-8 safe content.
//
// Compatible with web sockets and server sent events.
func Message(http *scopes.Http, message string) {
	Content(http, []byte(message))
}
