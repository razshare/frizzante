package negotiate

import (
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// SessionId tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func SessionId(http *scopes.Http) string {
	if http.SessionId != "" {
		return http.SessionId
	}
	var count uint
	var id string
	for _, cookie := range http.Request.CookiesNamed("session-id") {
		id = cookie.Value
		count++
	}
	if count > 0 {
		http.SessionId = id
		return id
	}
	// Create new session.
	ido, err := uuid.NewV4()
	if err != nil {
		logs.Errorf(
			http,
			"receive.SessionId: failed to create new session id: %v\n%s",
			err,
			stack.Trace(),
		)
		return ""
	}
	id = ido.String()
	send.Cookie(http, "session-id", id)
	http.SessionId = id
	return id
}
