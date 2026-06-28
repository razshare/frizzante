package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// BadRequest sends a message with status 400 Bad Request.
func BadRequest(http *scopes.Http, message string) {
	Status(http, http_.StatusBadRequest)
	Message(http, message)
}
