package send

import (
	"fmt"

	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// Redirectf redirects the request to a location with a status.
func Redirectf(http *scopes.Http, status int, location string, vars ...any) {
	Status(http, status)
	Header(http, "Location", fmt.Sprintf(location, vars...))
}
