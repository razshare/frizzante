package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/sessions"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func View() routes.Handler {
	return func(http *scopes.Http) {
		var session schema.Session
		sessions.Start(http, http.Queries, &session)
		defer func() {
			if err := http.Queries.ModifySessionById(http.Request.Context(), schema.ModifySessionByIdParams{
				ID: session.ID,
			}); err != nil {
				logs.Error(http, err)
			}
		}()
		context := http.Request.Context()
		var err error
		var todos []schema.Todo
		if todos, err = http.Queries.FindTodosBySessionId(context, session.ID); err != nil {
			session.Error = err.Error()
			return
		}
		send.View(http, views.View{Name: "Todos", Props: Props{
			Error: session.Error,
			Items: todos,
		}})
	}

}
