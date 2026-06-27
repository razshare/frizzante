package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

func Remove(queries *schema.Queries) routes.Handler {
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
			Id string `form:"id"`
		}
		if !receive.Form(client, &form) {
			session.Error = "could not parse form"
			return
		}
		context := client.Request.Context()
		if err := queries.RemoveTodosByIdAndSessionId(context, schema.RemoveTodosByIdAndSessionIdParams{
			ID:        form.Id,
			SessionID: session.ID,
		}); err != nil {
			session.Error = err.Error()
			return
		}
	}
}
