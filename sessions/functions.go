package sessions

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
)

// StartWithState starts a session with a given initial state.
//
// Deprecated: Use sessions.New followed by Session.Start.
func StartWithState[T any](con *connections.Connection, state T) *Session[T] {
	session := &Session[T]{
		Connection: con,
		State:      &state,
	}

	if !session.Exists() {
		session.Save()
	} else {
		session.Load()
	}

	return session
}

// Start starts a session with zero state.
//
// Deprecated: Use sessions.New followed by Session.Start.
func Start[T any](con *connections.Connection) *Session[T] {
	var state T
	return StartWithState(con, state)
}

// New creates a new session with a zero initial state.
func New[T any](con *connections.Connection, state T) *Session[T] {
	return &Session[T]{Connection: con, State: &state}
}

// Start loads the State if the current connection defines a session-id cookie.
//
// If the session-id cookie is missing it will create a new one and send it to the user.
func (session *Session[T]) Start() *Session[T] {
	if !session.Exists() {
		session.Save()
	} else {
		session.Load()
	}

	return session
}

// Id tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func (session *Session[T]) Id() string {
	if "" != session.Connection.SessionId {
		return session.Connection.SessionId
	}

	var sessionId string
	cookies := session.Connection.Request.CookiesNamed("session-id")
	count := 0

	for _, cookie := range cookies {
		sessionId = cookie.Value
		count++
	}

	if count > 0 {
		session.Connection.SessionId = sessionId
		return sessionId
	}

	// Create new session.
	id4, idError := uuid.NewV4()
	if idError != nil {
		session.Connection.Notifier.SendErrorAndTrace(idError, 1)
		session.Connection.SessionId = ""
		return ""
	}
	sessionId = id4.String()
	session.Connection.SendCookie("session-id", sessionId)
	session.Connection.SessionId = sessionId
	return sessionId
}

// Exists checks if the session exists into the archive.
func (session *Session[T]) Exists() bool {
	id := session.Id()
	has, hasError := session.Connection.SessionArchive.Has(id, globals.SessionKey)
	if hasError != nil {
		session.Connection.Notifier.SendErrorAndTrace(hasError, 1)
	}
	return has
}

// Save saves the session into the archive.
func (session *Session[T]) Save() *Session[T] {
	id := session.Id()
	readBytes, marshalError := json.Marshal(session.State)
	if marshalError != nil {
		session.Connection.Notifier.SendErrorAndTrace(marshalError, 1)
		return session
	}

	setError := session.Connection.SessionArchive.Set(id, globals.SessionKey, readBytes)
	if setError != nil {
		session.Connection.Notifier.SendErrorAndTrace(setError, 1)
	}
	return session
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func (session *Session[T]) Load() *Session[T] {
	id := session.Id()
	has, hasError := session.Connection.SessionArchive.Has(id, globals.SessionKey)
	if hasError != nil {
		session.Connection.Notifier.SendErrorAndTrace(hasError, 1)
		return session
	}
	if has {
		readBytes, getError := session.Connection.SessionArchive.Get(id, globals.SessionKey)
		if getError != nil {
			session.Connection.Notifier.SendErrorAndTrace(getError, 1)
			return session
		}
		unmarshalError := json.Unmarshal(readBytes, session.State)
		if unmarshalError != nil {
			session.Connection.Notifier.SendErrorAndTrace(unmarshalError, 1)
		}
	}
	return session
}

// Destroy removes the session from the archive.
func (session *Session[T]) Destroy() *Session[T] {
	id := session.Id()
	destroyError := session.Connection.SessionArchive.RemoveDomain(id)
	if destroyError != nil {
		session.Connection.Notifier.SendErrorAndTrace(destroyError, 1)
	}
	return session
}
