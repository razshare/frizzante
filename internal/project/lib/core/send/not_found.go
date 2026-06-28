package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// NotFound sends a message with status 404 Not Found.
func NotFound(http *scopes.Http, message string) {
	Status(http, http_.StatusNotFound)
	Message(http, message)
}
