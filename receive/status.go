package receive

import "github.com/razshare/frizzante/conn"

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(c *conn.Conn) <-chan struct{} {
	return c.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(c *conn.Conn) *bool {
	isAlive := true
	go func() {
		<-Cancellation(c)
		isAlive = false
	}()
	return &isAlive
}
