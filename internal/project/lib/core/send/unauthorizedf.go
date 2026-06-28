package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Unauthorizedf sends a message with status 401 Unauthorized.
func Unauthorizedf(http *scopes.Http, message string, vars ...any) {
	Status(http, http_.StatusUnauthorized)
	Messagef(http, message, vars...)
}
