package receive

import "github.com/razshare/frizzante/client"

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(c *client.Client) <-chan struct{} {
	return c.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(c *client.Client) *bool {
	isAlive := true
	go func() {
		<-Cancellation(c)
		isAlive = false
	}()
	return &isAlive
}
