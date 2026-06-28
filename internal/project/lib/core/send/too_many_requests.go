package send

import (
	http_ "net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// TooManyRequests sends a message with status 429 Too Many Requests.
func TooManyRequests(http *scopes.Http, message string) {
	Status(http, http_.StatusTooManyRequests)
	Message(http, message)
}
