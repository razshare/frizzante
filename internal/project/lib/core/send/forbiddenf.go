package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Forbiddenf sends a message with status 403 Forbidden.
func Forbiddenf(http *scopes.Http, message string, vars ...any) {
	Status(http, http_.StatusForbidden)
	Messagef(http, message, vars...)
}
