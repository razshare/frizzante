package todos

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/negotiate"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
)

func Remove(queries *schema.Queries) routes.Handler {
	return func(id uint64, request *http.Request, writer http.ResponseWriter) {
		var form struct {
			Id string `form:"id"`
		}
		context := request.Context()
		sessionId, _ := negotiate.SessionId(writer, request)
		session, _ := queries.FindSessionById(context, sessionId)
		_ = receive.Form(request, &form)
		_ = queries.RemoveTodosByIdAndSessionId(context, schema.RemoveTodosByIdAndSessionIdParams{
			ID:        form.Id,
			SessionID: session.ID,
		})
		_ = queries.ModifySessionById(request.Context(), schema.ModifySessionByIdParams{
			ID:    session.ID,
			Error: session.Error,
		})
		writer.Header().Set("Location", "/todos")
		writer.WriteHeader(302)
	}
}
