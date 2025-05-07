package frizzante

import (
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
)

var sessions = map[string]*Session{}

type Session struct {
	id       string
	get      func(key string, defaultValue any) (value any)
	set      func(key string, value any)
	unset    func(key string)
	validate func() (valid bool)
	destroy  func()
}

// SessionStart first tries to retrieve the client session, then,
// if the client is not associated with any existing session,
// it will automatically create a new empty session for said client.
//
// It always returns three functions, get, set and unset.
//
// Use get to retrieve a property from the session.
//
// Use set to create a new property or update an existing one to the session.
//
// Use unset to remove a property from the session.
func SessionStart(request *Request, response *Response) *Session {
	var sessionIdCookie *http.Cookie
	sessionIdCookies := request.httpRequest.CookiesNamed("session-id")
	sessionIdCookiesLen := len(sessionIdCookies)

	if 0 == sessionIdCookiesLen {
		uuidV4, sessionIdError := uuid.NewV4()

		if sessionIdError != nil {
			NotifierSendError(request.server.notifier, sessionIdError)
		}

		sessionId := uuidV4.String()
		freshSession := &Session{id: sessionId}
		request.server.sessionBuilder(freshSession)
		sessions[freshSession.id] = freshSession
		ResponseSendCookie(response, "session-id", freshSession.id)
		return freshSession
	}

	var sessionExists bool
	var session *Session

	for _, cookie := range sessionIdCookies {
		session, sessionExists = sessions[cookie.Value]
		if sessionExists {
			sessionIdCookie = cookie
			break
		}
	}

	if !sessionExists {
		uuidV4, sessionIdError := uuid.NewV4()
		if sessionIdError != nil {
			NotifierSendError(request.server.notifier, sessionIdError)
		}
		sessionId := uuidV4.String()
		freshSession := &Session{id: sessionId}
		request.server.sessionBuilder(freshSession)
		sessions[freshSession.id] = freshSession
		ResponseSendCookie(response, "session-id", sessionId)
		return freshSession
	}

	if !session.validate() {
		delete(sessions, sessionIdCookie.Value)
		session.destroy()
		return SessionStart(request, response)
	}

	ResponseSendCookie(response, "session-id", session.id)
	return session
}

// SessionGet gets a property from the session store.
func SessionGet[T any](self *Session, key string, defaultValue any) T {
	return self.get(key, defaultValue).(T)
}

// SessionSet sets a property in the session store.
func SessionSet[T any](self *Session, key string, value T) {
	self.set(key, value)
}

// SessionUnset unsets a property in the session store.
func SessionUnset(self *Session, key string) {
	self.unset(key)
}

// SessionValidate validates the session.
func SessionValidate(self *Session) bool {
	return self.validate()
}

// SessionDestroy destroys the session.
func SessionDestroy(self *Session) {
	self.destroy()
}

// SessionWithGetter sets the getter, which will retrieve a property from the session store.
func SessionWithGetter(self *Session, getter func(key string, defaultValue any) (value any)) {
	self.get = getter
}

// SessionWithSetter sets the setter, which will create or modify a property to the session store.
func SessionWithSetter(self *Session, setter func(key string, value any)) {
	self.set = setter
}

// SessionWithUnsetter sets the unsetter, which will remove a property from the session store.
func SessionWithUnsetter(self *Session, unsetter func(key string)) {
	self.unset = unsetter
}

// SessionWithValidator sets the validator, which will validate if a session is still alive.
func SessionWithValidator(self *Session, validator func() (valid bool)) {
	self.validate = validator
}

// SessionWithDestroyer sets the destroyer, which will destroy a session.
func SessionWithDestroyer(self *Session, destroyer func()) {
	self.destroy = destroyer
}
