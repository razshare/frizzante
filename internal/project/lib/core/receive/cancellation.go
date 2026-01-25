package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(client *clients.Client) <-chan struct{} {
	return client.Request.Context().Done()
}
