package fallback

import (
	"embed"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

func View(efs embed.FS) routes.Handler {
	return func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter) {
		if found, _ := send.RequestedFile(writer, request, efs, "/"); !found {
			writer.Header().Add("Location", "/welcome")
			writer.WriteHeader(http.StatusPermanentRedirect)
		}
	}
}
