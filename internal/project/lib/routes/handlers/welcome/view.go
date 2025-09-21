package welcome

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
)

type Props struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func View(client *client.Client) {
	send.View(client, view.View{
		Name: "Welcome",
		Props: Props{
			Message: "hello",
		},
	})
}
