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
	var id string
	cookies := c.Request.CookiesNamed("session-id")
	cookiesCount := 0

	for _, cookie := range cookies {
		id = cookie.Value
		cookiesCount++
	}

	if cookiesCount > 0 {
		return id
	}

	// Create new session.
	idObject, idObjectError := uuid.NewV4()
	if idObjectError != nil {
		c.Scope.Container.Config.ErrorLog.Println(idObjectError, stack.Trace())
		return ""
	}

	id = idObject.String()

	send.Cookie(c, "session-id", id)

	return id
}
