package todos

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func View(
	queries *schema.Queries,
	render renders.Render,
) routes.Handler {
	return func(scope routes.Scope, request *http.Request, writer http.ResponseWriter) {
		context := request.Context()
		session := scope["session"].(schema.Session)
		todos, _ := queries.FindTodosBySessionId(context, session.ID)
		_ = queries.ModifySessionById(context, schema.ModifySessionByIdParams{
			ID: session.ID,
		})
		_ = send.View(writer, request, render, views.View{
			Name: "Todos",
			Props: Props{
				Error: session.Error,
				Items: todos,
			},
		})
	}

}
