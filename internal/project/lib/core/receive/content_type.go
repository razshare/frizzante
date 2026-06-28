package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// ContentType reads the Content-Type header field and returns the value.
func ContentType(http *scopes.Http) string {
	return http.Request.Header.Get("Content-Type")
}
