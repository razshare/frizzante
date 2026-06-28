package send

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Navigatef redirects the request to a location with status 302.
func Navigatef(http *scopes.Http, format string, vars ...any) {
	Redirect(http, fmt.Sprintf(format, vars...), 302)
	Message(http, "")
}
