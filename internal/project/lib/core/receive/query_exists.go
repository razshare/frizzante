package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// QueryExists checks if a query field exists.
func QueryExists(http *scopes.Http, key string) bool {
	return http.Request.URL.Query().Has(key)
}
