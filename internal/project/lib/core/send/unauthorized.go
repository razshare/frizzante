package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Unauthorized sends a message with status 401 Unauthorized.
func Unauthorized(http *scopes.Http, message string) {
	Status(http, http_.StatusUnauthorized)
	Message(http, message)
}
