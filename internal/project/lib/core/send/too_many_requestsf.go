package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// TooManyRequestsf sends a message with status 403 Forbidden.
func TooManyRequestsf(client *clients.Client, message string, vars ...any) {
	Status(client, http.StatusTooManyRequests)
	Messagef(client, message, vars...)
}
