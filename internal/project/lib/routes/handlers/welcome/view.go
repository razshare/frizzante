package welcome

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
)

func View(client *client.Client) {
	send.View(client, view.View{Name: "Welcome"})
}
