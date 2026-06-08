package receive

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/databases"
	"github.com/razshare/frizzante/internal/project/lib/core/databases/schema"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

func Session(client *clients.Client, session *schema.Session) bool {
	id := SessionId(client)
	context := client.Request.Context()
	var err error
	if *session, err = databases.Queries.FindSessionById(context, id); err != nil {
		logs.Errorf(client, "something went wrong while retrieving session %s: %v", id, err)
		logs.Infof(client, "attempting to add session %s...", id)
		if err = databases.Queries.AddSessionWithId(context, id); err != nil {
			logs.Errorf(client, "attempt to add session %s failed: %v\n%s", id, err, stack.Trace())
			return false
		}
		logs.Infof(client, "session %s created", id)
		session.ID = id
	}
	return true
}
