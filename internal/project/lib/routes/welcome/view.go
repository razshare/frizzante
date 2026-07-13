package welcome

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func View(render renders.Render) routes.Handler {
	return func(id uint64, request *http.Request, writer http.ResponseWriter) {
		_ = send.View(writer, request, render, views.View{
			Name: "Welcome",
		})
	}
}
