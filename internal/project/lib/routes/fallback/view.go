package fallback

import (
	"embed"
	"net/http"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/send"
)

func View(efs embed.FS) routes.Handler {
	return func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter) {
		if found, _ := send.RequestedFile(writer, request, efs, "/"); !found {
			send.ToLocation(writer, "/welcome")
		}
	}
}
