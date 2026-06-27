package todos

import (
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

func Add(queries *schema.Queries) routes.Handler {
	return func(client *clients.Client) {
		var session schema.Session
		receive.Session(client, queries, &session)
		defer send.Navigate(client, "/todos")
		defer func() {
			if err := queries.ModifySessionById(client.Request.Context(), schema.ModifySessionByIdParams{
				ID:    session.ID,
				Error: session.Error,
			}); err != nil {
				logs.Error(client, err)
			}
		}()
		var form struct {
			Description string `form:"description"`
		}
		if !receive.Form(client, &form) {
			session.Error = "could not parse form"
			return
		}
		if form.Description == "" {
			session.Error = "description cannot be empty"
			return
		}
		ido, err := uuid.NewV4()
		if err != nil {
			session.Error = err.Error()
			return
		}
		id := ido.String()
		context := client.Request.Context()
		if err = queries.AddTodoWithIdAndSessionId(context, schema.AddTodoWithIdAndSessionIdParams{
			ID:          id,
			SessionID:   session.ID,
			Description: form.Description,
		}); err != nil {
			session.Error = err.Error()
			return
		}
	}
}
