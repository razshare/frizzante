package todos

import (
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/keys"
	schema2 "github.com/razshare/frizzante/internal/project/lib/schema"
)

func Add(queries *schema2.Queries) routes.Handler {
	return func(scope scopes.Scope, request *http.Request, writer http.ResponseWriter) {
		var form struct {
			Description string `form:"description"`
		}
		context := request.Context()
		session := scope[keys.Session].(schema2.Session)
		_ = receive.Form(request, &form)
		if form.Description == "" {
			session.Error = "description cannot be empty"
			return
		}
		todoId, _ := uuid.NewV4()
		_ = queries.AddTodoWithIdAndSessionId(context, schema2.AddTodoWithIdAndSessionIdParams{
			ID:          todoId.String(),
			SessionID:   session.ID,
			Description: form.Description,
		})
		_ = queries.ModifySessionById(request.Context(), schema2.ModifySessionByIdParams{
			ID:    session.ID,
			Error: session.Error,
		})
		writer.Header().Set("Location", "/todos")
		writer.WriteHeader(302)
	}
}
