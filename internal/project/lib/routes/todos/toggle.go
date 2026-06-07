package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/databases"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/sqlc"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

func Toggle(client *clients.Client) {
	var session sqlc.Session
	defer send.Navigate(client, "/todos")
	defer func() {
		if err := databases.Queries.ModifySessionById(client.Request.Context(), sqlc.ModifySessionByIdParams{
			ID:    session.ID,
			Error: session.Error,
		}); err != nil {
			logs.Error(client, err)
		}
	}()
	receive.Session(client, &session)
	var form struct {
		Id    string `form:"id"`
		Value int64  `form:"value"`
	}
	if !receive.Form(client, &form) {
		session.Error = "could not parse form"
		return
	}
	if err := databases.Queries.ToggleTodosByIdAndSessionId(client.Request.Context(), sqlc.ToggleTodosByIdAndSessionIdParams{
		ID:        form.Id,
		SessionID: session.ID,
		Checked:   form.Value,
	}); err != nil {
		session.Error = err.Error()
		return
	}
}
