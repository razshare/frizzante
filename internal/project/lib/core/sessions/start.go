package sessions

import (
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/negotiate"
)

func Start(request *http.Request, writer http.ResponseWriter, queries *schema.Queries, session *schema.Session) (err error) {
	id, _ := negotiate.SessionId(request, writer)
	context := request.Context()
	if *session, err = queries.FindSessionById(context, id); err != nil {
		if err = queries.AddSessionWithId(context, id); err != nil {
			return
		}
		session.ID = id
	}
	return
}
