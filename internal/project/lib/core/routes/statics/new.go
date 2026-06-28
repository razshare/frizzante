package statics

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/servers"
)

// NewRouteHandler creates a route handler that lists all static routes of a given server.
func NewRouteHandler(server *servers.Server) routes.Handler {
	return func(request *http.Request, writer http.ResponseWriter) {
		if accepts := request.Header.Get("Accept"); accepts != "" && accepts != "application/json" {
			writer.WriteHeader(http.StatusBadRequest)
			_ = send.Message(writer, fmt.Sprintf("only application/json can be produced; requested %s", accepts))
			return
		}
		statics := make([]string, 0)
		for _, route := range server.Routes {
			if parts := strings.SplitN(route.Pattern, " ", 2); len(parts) >= 2 && parts[0] == "GET" {
				statics = append(statics, parts[1])
			}
		}
		_ = send.Json(writer, statics)
	}
}
