package todos

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/routes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/sessions"
	"github.com/razshare/frizzante/internal/project/lib/core/views"
	"github.com/razshare/frizzante/internal/project/lib/core/views/renders"
)

func View(
	queries *schema.Queries,
	render renders.Render,
) routes.Handler {
	return func(request *http.Request, writer http.ResponseWriter) {
		var session schema.Session
		_ = sessions.Start(request, writer, queries, &session)
		todos, _ := queries.FindTodosBySessionId(request.Context(), session.ID)
		_ = queries.ModifySessionById(request.Context(), schema.ModifySessionByIdParams{
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
