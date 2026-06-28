package welcome

import (
	"embed"
	"log"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func View(
	render renders.Render,
	efs embed.FS,
	logerr *log.Logger,
	loginf *log.Logger,
) routes.Handler {
	return func(request *http.Request, writer http.ResponseWriter) {
		_ = send.View(writer, request, render, efs, logerr, loginf, views.View{
			Name: "Welcome",
		})
	}
}
