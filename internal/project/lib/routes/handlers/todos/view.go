package todos

import (
	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
	_session "github.com/razshare/frizzante/internal/project/lib/session"
)

func View(client *_client.Client) {
	session := _session.Start(receive.SessionId(client))
	send.View(client, view.View{
		Name: "Todos",
		Props: Props{
			Items: session.Todos,
			Error: receive.Query(client, "error"),
		},
	})
}
