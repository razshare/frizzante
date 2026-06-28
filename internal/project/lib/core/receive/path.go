package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Path reads a path value field and returns it.
func Path(http *scopes.Http, key string) string {
	return http.Request.PathValue(key)
}
