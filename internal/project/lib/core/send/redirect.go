package send

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Redirect redirects the request to a location with a status.
func Redirect(client *clients.Client, location string, status int) {
	Status(client, status)
	Header(client, "Location", location)
}
