package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Error sends a message with status 500 Internal Server Error.
func Error(http *scopes.Http, err error) {
	type ServerError struct {
		Error string
	}
	Status(http, http_.StatusInternalServerError)
	Message(http, err.Error())
}
