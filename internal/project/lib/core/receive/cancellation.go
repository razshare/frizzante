package receive

import _client "github.com/razshare/frizzante/internal/project/lib/core/client"

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(client *_client.Client) <-chan struct{} {
	return client.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(client *_client.Client) *bool {
	alive := true
	go func() {
		<-Cancellation(client)
		alive = false
	}()
	return &alive
}
