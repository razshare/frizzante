package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Forbidden sends a message with status 403 Forbidden.
func Forbidden(http *scopes.Http, message string) {
	Status(http, http_.StatusForbidden)
	Message(http, message)
}
