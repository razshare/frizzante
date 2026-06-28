package todos

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/sessions"
)

func Remove(queries *schema.Queries) routes.Handler {
	return func(request *http.Request, writer http.ResponseWriter) {
		var session schema.Session
		_ = sessions.Start(request, writer, queries, &session)
		var form struct {
			Id string `form:"id"`
		}
		_ = receive.Form(request, &form)
		context := request.Context()
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
