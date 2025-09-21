package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
	"github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func View(client *client.Client) {
	state := session.Start(receive.SessionId(client))
	send.View(client, view.View{
		Name: "Todos",
		Props: Props{
			Todos: state.Todos,
			Error: receive.Query(client, "error"),
		},
	})
}
