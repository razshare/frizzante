package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(http *scopes.Http) *bool {
	alive := true
	go func() {
		<-Cancellation(http)
		alive = false
	}()
	return &alive
}
