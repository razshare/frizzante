package todos

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/receive"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	_view "github.com/razshare/frizzante/internal/project/lib/core/view"
	"github.com/razshare/frizzante/internal/project/lib/session/memory"
)

func Remove(c *client.Client) {
	s := session.Start(receive.SessionId(c))

	l := int64(len(s.Todos))
	if 0 == l {
		// No index found, ignore the request.
		send.Navigate(c, "/todos")
		return
	}

	is := receive.Query(c, "index")
	if is == "" {
		// No index found, ignore the request.
		send.Navigate(c, "/todos")
		return
	}

	i, e := strconv.ParseInt(is, 10, 64)
	if nil != e {
		send.View(c, _view.View{
			Name: "Todos",
			Props: map[string]any{
				"todos": s.Todos,
				"error": e.Error(),
			},
		})
		return
	}
	if i >= l {
		// Index is out of bounds, ignore the request.
		send.Navigate(c, "/todos")
		return
	}

	s.Todos = append(
		s.Todos[:i],
		s.Todos[i+1:]...,
	)

	send.Navigate(c, "/todos")
}
