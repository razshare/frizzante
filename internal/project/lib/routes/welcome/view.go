package welcome

import (
	"net/http"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/send"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/views"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/views/renders"
)

func View(render renders.Render) routes.Handler {
	return func(_ scopes.Scope, request *http.Request, writer http.ResponseWriter) {
		_ = send.View(writer, request, render, views.View{
			Name: "Welcome",
		})
	}
}
