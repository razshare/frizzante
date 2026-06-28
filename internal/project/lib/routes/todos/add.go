package todos

import (
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/sessions"
)

func Add(queries *schema.Queries) routes.Handler {
	return func(request *http.Request, writer http.ResponseWriter) {
		var session schema.Session
		_ = sessions.Start(request, writer, queries, &session)
		var form struct {
			Description string `form:"description"`
		}
		_ = receive.Form(request, &form)
		if form.Description == "" {
			session.Error = "description cannot be empty"
			return
		}
		todoId, _ := uuid.NewV4()
		context := request.Context()
		_ = queries.AddTodoWithIdAndSessionId(context, schema.AddTodoWithIdAndSessionIdParams{
			ID:          todoId.String(),
			SessionID:   session.ID,
			Description: form.Description,
		})
		_ = queries.ModifySessionById(request.Context(), schema.ModifySessionByIdParams{
			ID:    session.ID,
			Error: session.Error,
		})
		writer.Header().Set("Location", "/todos")
		writer.WriteHeader(302)
	}
}
