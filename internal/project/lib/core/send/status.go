package send

import (
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Status sets the status code.
//
// This will lock the status, which makes it
// so that the next time you invoke this
// function it will fail with an error.
//
// All errors are sent to the server notifier.
func Status(http *scopes.Http, status int) {
	if http.Locked {
		logs.Errorf(
			http,
			"send.Status: status is locked, cannot set status\n%s",
			stack.Trace(),
		)
		return
	}
	http.Status = status
}
