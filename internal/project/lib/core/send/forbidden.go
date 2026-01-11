package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// Forbidden sends a message with status 403 Forbidden.
func Forbidden(client *clients.Client, message string) {
	Status(client, http.StatusForbidden)
	Message(client, message)
}
