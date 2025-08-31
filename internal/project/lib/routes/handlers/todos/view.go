package todos

import (
	"main/lib/session/memory"

	"main/lib/core/client"
	"main/lib/core/receive"
	"main/lib/core/send"
	"main/lib/core/view"
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
