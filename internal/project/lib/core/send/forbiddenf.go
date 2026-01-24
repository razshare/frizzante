package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// Forbiddenf sends a message with status 403 Forbidden.
func Forbiddenf(client *clients.Client, message string, vars ...any) {
	Status(client, http.StatusForbidden)
	Messagef(client, message, vars...)
}
