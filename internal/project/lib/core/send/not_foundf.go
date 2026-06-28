package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// NotFoundf sends a message with status 404 Not Found.
func NotFoundf(http *scopes.Http, message string, vars ...any) {
	Status(http, http_.StatusNotFound)
	Messagef(http, message, vars...)
}
