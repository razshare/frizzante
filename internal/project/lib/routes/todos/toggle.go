package todos

import (
	"net/http"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/v2/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/v2/internal/project/lib/keys"
	schema2 "github.com/razshare/frizzante/v2/internal/project/lib/schema"
)

func Toggle(queries *schema2.Queries) routes.Handler {
	return func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter) {
		var form struct {
			Id    string `form:"id"`
			Value int64  `form:"value"`
		}
		context := request.Context()
		session := scope[keys.Session].(schema2.Session)
		_ = receive.Form(request, &form)
		_ = queries.ToggleTodosByIdAndSessionId(context, schema2.ToggleTodosByIdAndSessionIdParams{
			ID:        form.Id,
			SessionID: session.ID,
			Checked:   form.Value,
		})
		_ = queries.ModifySessionById(context, schema2.ModifySessionByIdParams{
			ID:    session.ID,
			Error: session.Error,
		})
		writer.Header().Set("Location", "/todos")
		writer.WriteHeader(302)
	}
}
