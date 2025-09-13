package welcome

import (
	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/view"
)

func View(c *client.Client) {
	send.View(c, view.View{Name: "Welcome"})
}
