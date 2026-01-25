package send

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Navigate redirects the request to a location with status 302.
func Navigate(client *clients.Client, location string) {
	Redirect(client, location, 302)
	Message(client, "")
}
