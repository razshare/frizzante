package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Header reads a header field and returns the value.
func Header(http *scopes.Http, key string) string {
	return http.Request.Header.Get(key)
}
