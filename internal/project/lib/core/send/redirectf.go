package send

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// Redirectf redirects the request to a location with a status.
func Redirectf(client *clients.Client, status int, location string, vars ...any) {
	Status(client, status)
	Header(client, "Location", fmt.Sprintf(location, vars...))
}
