package todos

import (
	"main/lib/session/memory"

	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/receive"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/view"
)

func View(c *client.Client) {
	s := session.Start(receive.SessionId(c))
	send.View(c, view.View{
		Name: "Todos",
		Props: map[string]any{
			"todos": s.Todos,
		},
	})
}
