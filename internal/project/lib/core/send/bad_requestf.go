package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// BadRequestf sends a message with status 400 Bad Request.
func BadRequestf(http *scopes.Http, message string, vars ...any) {
	Status(http, http_.StatusBadRequest)
	Messagef(http, message, vars...)
}
