package receive

import (
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/stack"
)

// SessionId tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func SessionId(c *client.Client) string {
	if c.SessionId != "" {
		return c.SessionId
	}

	var j uint
	var id string

	for _, ck := range c.Request.CookiesNamed("session-id") {
		id = ck.Value
		j++
	}

	if j > 0 {
		c.SessionId = id
		return id
	}

	// Create new session.
	ido, err := uuid.NewV4()
	if err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return ""
	}

	id = ido.String()

	send.Cookie(c, "session-id", id)

	c.SessionId = id

	return id
}
