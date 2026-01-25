package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Accept reads if the Accept header entries and returns the values.
func Accept(client *clients.Client) string {
	return client.Request.Header.Get("Accept")
}
