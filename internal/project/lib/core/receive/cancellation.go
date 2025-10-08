package receive

import "github.com/razshare/frizzante/internal/project/lib/core/clients"

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(client *clients.Client) <-chan struct{} {
	return client.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(client *clients.Client) *bool {
	alive := true
	go func() {
		<-Cancellation(client)
		alive = false
	}()
	return &alive
}
