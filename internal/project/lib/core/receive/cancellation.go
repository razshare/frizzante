package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(http *scopes.Http) <-chan struct{} {
	return http.Request.Context().Done()
}
