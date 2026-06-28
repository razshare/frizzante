package send

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Messagef sends utf-8 safe content using a format.
//
// Compatible with web sockets and server sent events.
func Messagef(http *scopes.Http, format string, vars ...any) {
	Content(http, []byte(fmt.Sprintf(format, vars...)))
}
