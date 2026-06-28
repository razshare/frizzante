package fallback

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
)

func View(
	render renders.Render,
	efs embed.FS,
	logerr *log.Logger,
	loginf *log.Logger,
) routes.Handler {
	view := welcome.View(render, efs, logerr, loginf)
	return func(request *http.Request, writer http.ResponseWriter) {
		if found, _ := send.RequestedFile(writer, request, efs); !found {
			view(request, writer)
		}
	}
}
