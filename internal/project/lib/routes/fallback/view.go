package fallback

import (
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/routes/welcome"
)

func View() routes.Handler {
	return func(http *scopes.Http) {
		if !send.RequestedFile(http) {
			welcome.View()(http)
		}
	}
}
