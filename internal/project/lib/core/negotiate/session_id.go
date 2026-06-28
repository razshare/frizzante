package negotiate

import (
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/internal/project/lib/core/send"
)

// SessionId tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func SessionId(request *http.Request, writer http.ResponseWriter) (value string, err error) {
	var cookie *http.Cookie
	if cookie, err = request.Cookie("session-id"); err != nil {
		return
	}
	if value = cookie.Value; value != "" {
		return
	}
	var ido *uuid.UUID
	if ido, err = uuid.NewV4(); err != nil {
		return
	}
	value = ido.String()
	request.AddCookie(&http.Cookie{Name: "session-id", Value: value})
	send.Cookie(writer, "session-id", value)
	return
}
