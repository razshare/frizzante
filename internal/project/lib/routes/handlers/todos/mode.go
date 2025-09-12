package todos

import (
	"main/lib/core/client"
	"main/lib/core/receive"
	"main/lib/core/send"
	session "main/lib/session/memory"
)

func Mode(c *client.Client) {
	defer send.Navigate(c, "/todos")

	var val string
	if val = receive.Query(c, "value"); val == "" {
		return
	}

	s := session.Start(receive.SessionId(c))

	switch val {
	case "remove":
		s.Mode = session.ModeRemove
	case "add":
		s.Mode = session.ModeAdd
	default:
		s.Mode = session.ModeToggle
	}
}
