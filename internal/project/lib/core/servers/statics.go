package servers

import (
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

func Statics(server *Server, pattern string) routes.Route {
	return routes.Route{Pattern: pattern, Handler: func(client *clients.Client) {
		if accepts := receive.Accept(client); accepts != "" && accepts != "application/json" {
			send.BadRequestf(client, "only application/json can be produced; requested %s", accepts)
			return
		}
		statics := make([]string, 0)
		for _, route := range server.Routes {
			if parts := strings.SplitN(route.Pattern, " ", 2); len(parts) >= 2 && parts[0] == "GET" {
				statics = append(statics, parts[1])
			}
		}
		send.Json(client, statics)
	}}
}
