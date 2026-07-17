package statics

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

// NewRouteHandler creates a route handler that lists all static routes of a given server.
func NewRouteHandler(appRoutes []routes.Route) routes.Handler {
	return func(_ scopes.Scope, request *http.Request, writer http.ResponseWriter) {
		if accepts := request.Header.Get("Accept"); accepts != "application/json" {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(writer, "only application/json can be produced; requested %s", accepts)
			return
		}
		statics := make([]string, 0)
		for _, route := range appRoutes {
			if parts := strings.SplitN(route.Pattern, " ", 2); len(parts) >= 2 && parts[0] == "GET" {
				statics = append(statics, parts[1])
			}
		}
		_ = send.Json(writer, statics)
	}
}
