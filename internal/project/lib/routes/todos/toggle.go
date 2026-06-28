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

func Toggle() routes.Handler {
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
			Id    string `form:"id"`
			Value int64  `form:"value"`
		}
		if !receive.Form(http, &form) {
			session.Error = "could not parse form"
			return
		}
		if err := http.Queries.ToggleTodosByIdAndSessionId(http.Request.Context(), schema.ToggleTodosByIdAndSessionIdParams{
			ID:        form.Id,
			SessionID: session.ID,
			Checked:   form.Value,
		}); err != nil {
			session.Error = err.Error()
			return
		}
	}
}
