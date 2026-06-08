package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/databases"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
)

func View(client *clients.Client) {
	var session schema.Session
	defer func() {
		if err := databases.Queries.ModifySessionById(client.Request.Context(), schema.ModifySessionByIdParams{
			ID: session.ID,
		}); err != nil {
			logs.Error(client, err)
		}
	}()
	receive.Session(client, &session)
	context := client.Request.Context()
	var err error
	var todos []schema.Todo
	if todos, err = databases.Queries.FindTodosBySessionId(context, session.ID); err != nil {
		session.Error = err.Error()
		return
	}
	send.View(client, views.View{Name: "Todos", Props: Props{
		Error: session.Error,
		Items: todos,
	}})
}
