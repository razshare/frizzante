package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/sessions"
)

func Remove() routes.Handler {
	return func(http *scopes.Http) {
		var session schema.Session
		sessions.Start(http, http.Queries, &session)
		defer send.Navigate(http, "/todos")
		defer func() {
			if err := http.Queries.ModifySessionById(http.Request.Context(), schema.ModifySessionByIdParams{
				ID:    session.ID,
				Error: session.Error,
			}); err != nil {
				logs.Error(http, err)
			}
		}()
		var form struct {
			Id string `form:"id"`
		}
		if !receive.Form(http, &form) {
			session.Error = "could not parse form"
			return
		}
		context := http.Request.Context()
		if err := http.Queries.RemoveTodosByIdAndSessionId(context, schema.RemoveTodosByIdAndSessionIdParams{
			ID:        form.Id,
			SessionID: session.ID,
		}); err != nil {
			session.Error = err.Error()
			return
		}
	}
}
