package todos

import (
	"main/lib/session/memory"
	"strconv"

	"main/lib/core/client"
	"main/lib/core/receive"
	"main/lib/core/send"
	"main/lib/core/view"
)

func Uncheck(c *client.Client) {
	s := session.Start(receive.SessionId(c))

	is := receive.Query(c, "index")
	if is == "" {
		// No index found, ignore the request.
		send.Navigate(c, "/todos")
		return
	}

	i, e := strconv.ParseInt(is, 10, 64)
	if nil != e {
		send.View(c, view.View{
			Name: "Todos",
			Props: map[string]any{
				"todos": s.Todos,
				"error": e.Error(),
			},
		})
		return
	}

	l := int64(len(s.Todos))

	if i >= l {
		// Index is out of bounds, ignore the request.
		send.Navigate(c, "/todos")
		return
	}

	s.Todos[i].Checked = false

	send.Navigate(c, "/todos")
}
