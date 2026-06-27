package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func View(queries *schema.Queries) routes.Handler {
	return func(client *clients.Client) {
		var session schema.Session
		receive.Session(client, queries, &session)
		defer func() {
			if err := queries.ModifySessionById(client.Request.Context(), schema.ModifySessionByIdParams{
				ID: session.ID,
			}); err != nil {
				logs.Error(client, err)
			}
		}()
		context := client.Request.Context()
		var err error
		var todos []schema.Todo
		if todos, err = queries.FindTodosBySessionId(context, session.ID); err != nil {
			session.Error = err.Error()
			return
		}
		send.View(client, views.View{Name: "Todos", Props: Props{
			Error: session.Error,
			Items: todos,
		}})
	}

}
