package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// Unauthorized sends a message with status 401 Unauthorized.
func Unauthorized(client *clients.Client, message string) {
	Status(client, http.StatusUnauthorized)
	Message(client, message)
}
