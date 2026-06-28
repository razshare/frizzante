package fallback

import (
	"embed"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
)

func View(
	render renders.Render,
	efs embed.FS,
) routes.Handler {
	view := welcome.View(render)
	return func(request *http.Request, writer http.ResponseWriter) {
		if found, _ := send.RequestedFile(writer, request, efs); !found {
			view(request, writer)
		}
	}
}
