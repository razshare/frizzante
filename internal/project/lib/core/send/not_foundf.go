package send

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// NotFoundf sends a message with status 404 Not Found.
func NotFoundf(client *clients.Client, message string, vars ...any) {
	Status(client, http.StatusNotFound)
	Messagef(client, message, vars...)
}
