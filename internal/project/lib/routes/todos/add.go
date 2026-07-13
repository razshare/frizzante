package todos

import (
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/negotiate"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
)

func Add(queries *schema.Queries) routes.Handler {
	return func(id uint64, request *http.Request, writer http.ResponseWriter) {
		var form struct {
			Description string `form:"description"`
		}
		context := request.Context()
		sessionId, _ := negotiate.SessionId(writer, request)
		session, _ := queries.FindSessionById(context, sessionId)
		_ = receive.Form(request, &form)
		if form.Description == "" {
			session.Error = "description cannot be empty"
			return
		}
		todoId, _ := uuid.NewV4()
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
