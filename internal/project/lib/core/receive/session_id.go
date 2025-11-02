package receive

import (
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// SessionId tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func SessionId(client *clients.Client) string {
	if client.SessionId != "" {
		return client.SessionId
	}

	var count uint
	var id string

	for _, cookie := range client.Request.CookiesNamed("session-id") {
		id = cookie.Value
		count++
	}

	if count > 0 {
		client.SessionId = id
		return id
	}

	// Create new session.
	ido, err := uuid.NewV4()
	if err != nil {
		client.Options.ErrorLog.Println(err, stack.Trace())
		return ""
	}

	id = ido.String()

	send.Cookie(client, "session-id", id)

	client.SessionId = id

	return id
}
