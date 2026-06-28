package todos

import (
	"errors"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/sessions"
)

func Add() routes.Handler {
	return func(http *scopes.Http) {
		var session schema.Session
		if !sessions.Start(http, http.Queries, &session) {
			send.Error(http, errors.New("could not start session"))
			return
		}
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
			Description string `form:"description"`
		}
		if !receive.Form(http, &form) {
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
		context := http.Request.Context()
		if err = http.Queries.AddTodoWithIdAndSessionId(context, schema.AddTodoWithIdAndSessionIdParams{
			ID:          id,
			SessionID:   session.ID,
			Description: form.Description,
		}); err != nil {
			session.Error = err.Error()
			return
		}
	}
}
