package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// TooManyRequestsf sends a message with status 403 Forbidden.
func TooManyRequestsf(http *scopes.Http, message string, vars ...any) {
	Status(http, http_.StatusTooManyRequests)
	Messagef(http, message, vars...)
}
