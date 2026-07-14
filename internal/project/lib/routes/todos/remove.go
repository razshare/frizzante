package todos

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/keys"
)

func Remove(queries *schema.Queries) routes.Handler {
	return func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter) {
		var form struct {
			Id string `form:"id"`
		}
		context := request.Context()
		session := scope[keys.Session].(schema.Session)
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
