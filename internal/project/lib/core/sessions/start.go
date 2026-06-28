package sessions

import (
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/negotiate"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func Start(http *scopes.Http, queries *schema.Queries, session *schema.Session) bool {
	id := negotiate.SessionId(http)
	context := http.Request.Context()
	var err error
	if *session, err = queries.FindSessionById(context, id); err != nil {
		logs.Errorf(http, "something went wrong while retrieving session %s: %v", id, err)
		logs.Infof(http, "attempting to add session %s...", id)
		if err = queries.AddSessionWithId(context, id); err != nil {
			logs.Errorf(http, "attempt to add session %s failed: %v\n%s", id, err, stack.Trace())
			return false
		}
		logs.Infof(http, "session %s created", id)
		session.ID = id
	}
	return true
}
