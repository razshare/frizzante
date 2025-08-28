package welcome

import (
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/view"
)

func View(c *client.Client) {
	send.View(c, view.View{Name: "Welcome"})
}
