package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// BadRequestf sends a message with status 400 Bad Request.
func BadRequestf(client *clients.Client, message string, vars ...any) {
	Status(client, http.StatusBadRequest)
	Messagef(client, message, vars...)
}
