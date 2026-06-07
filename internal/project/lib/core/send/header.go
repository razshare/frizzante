package send

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
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
		logs.Errorf(client, "header is locked\n%s", stack.Trace())
		return
	}
	client.Writer.Header().Set(key, value)
}
