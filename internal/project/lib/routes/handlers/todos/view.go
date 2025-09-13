package todos

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
	"github.com/razshare/frizzante/internal/project/lib/session/memory"
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
