package receive

import "github.com/razshare/frizzante/internal/project/lib/core/scopes"

// Accept reads if the Accept header entries and returns the values.
func Accept(http *scopes.Http) string {
	return http.Request.Header.Get("Accept")
}
