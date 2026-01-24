package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// Unauthorizedf sends a message with status 401 Unauthorized.
func Unauthorizedf(client *clients.Client, message string, vars ...any) {
	Status(client, http.StatusUnauthorized)
	Messagef(client, message, vars...)
}
