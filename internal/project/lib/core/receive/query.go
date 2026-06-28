package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Query reads a query field and returns the value.
func Query(http *scopes.Http, key string) string {
	return http.Request.URL.Query().Get(key)
}
