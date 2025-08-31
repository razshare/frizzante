package receive

import "main/lib/core/client"

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(client *client.Client) <-chan struct{} {
	return client.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(client *client.Client) *bool {
	alive := true
	go func() {
		<-Cancellation(client)
		alive = false
	}()
	return &alive
}
