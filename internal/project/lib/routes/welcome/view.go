package welcome

import (
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func View() routes.Handler {
	return func(http *scopes.Http) {
		send.View(http, views.View{Name: "Welcome"})
	}
}
