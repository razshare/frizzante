package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// TooManyRequests sends a message with status 429 Too Many Requests.
func TooManyRequests(client *clients.Client, message string) {
	Status(client, http.StatusTooManyRequests)
	Message(client, message)
}
