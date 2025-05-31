package frizzante

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
)

type SessionOperator[T any] interface {
	Id() string
	Load()
	Exists() bool
	Save()
	Destroy()
	Start(func(state T, save func(T)))
}

type SessionHandler[T any] = func(state *T)
type SessionGuard[T any] = func(connection *Connection, state *T, pass func())

type Session[T any] struct {
	connection *Connection
	archive    Archive
	guards     []SessionGuard[T]
	state      T
}

// NewSession creates a new session.
func NewSession[T any](connection *Connection) *Session[T] {
	archive := NewDiskArchive().WithName(".sessions")
	return &Session[T]{
		connection: connection,
		archive:    archive,
	}
}

// WithState sets the initial state.
func (session *Session[T]) WithState(state T) *Session[T] {
	session.state = state
	return session
}

// WithArchive sets the archive.
func (session *Session[T]) WithArchive(archive Archive) *Session[T] {
	session.archive = archive
	return session
}

// WithGuards sets the guards.
func (session *Session[T]) WithGuards(guards []SessionGuard[T]) *Session[T] {
	session.guards = guards
	return session
}

// Id tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func (session *Session[T]) Id() string {
	if "" != session.connection.sessionId {
		return session.connection.sessionId
	}

	var sessionId string
	cookies := session.connection.request.CookiesNamed("session-id")
	count := 0

	for _, cookie := range cookies {
		sessionId = cookie.Value
		count++
	}

	if count > 0 {
		session.connection.sessionId = sessionId
		return sessionId
	}

	// Create new session.
	id4, idError := uuid.NewV4()
	if idError != nil {
		session.connection.server.notifier.SendError(idError)
		session.connection.sessionId = ""
		return ""
	}
	sessionId = id4.String()
	session.connection.SendCookie("session-id", sessionId)
	session.connection.sessionId = sessionId
	return sessionId
}

// Load loads the session from the archive.
func (session *Session[T]) Load() {
	id := session.Id()
	has, hasError := session.archive.Has(id, SessionKey)
	if hasError != nil {
		session.connection.server.notifier.SendError(hasError)
		return
	}
	if has {
		readBytes, getError := session.archive.Get(id, SessionKey)
		if getError != nil {
			session.connection.server.notifier.SendError(getError)
			return
		}
		unmarshalError := json.Unmarshal(readBytes, &session.state)
		if unmarshalError != nil {
			session.connection.server.notifier.SendError(unmarshalError)
		}
	}
}

// Exists checks if the session exists into the archive.
func (session *Session[T]) Exists() bool {
	id := session.Id()
	has, hasError := session.archive.Has(id, SessionKey)
	if hasError != nil {
		session.connection.server.notifier.SendError(hasError)
	}
	return has
}

// Save saves the session into the archive.
func (session *Session[T]) Save() {
	id := session.Id()
	readBytes, marshalError := json.Marshal(&session.state)
	if marshalError != nil {
		session.connection.server.notifier.SendError(marshalError)
		return
	}
	setError := session.archive.Set(id, SessionKey, readBytes)
	if setError != nil {
		session.connection.server.notifier.SendError(setError)
	}
}

// Destroy removes the session from the archive.
func (session *Session[T]) Destroy() {
	id := session.Id()
	destroyError := session.archive.RemoveDomain(id)
	if destroyError != nil {
		session.connection.server.notifier.SendError(destroyError)
	}
}

// Start checks if the session given by the user already exists in the archive.
// If it does exist, it loads the state, otherwise it saves the state into the archive.
//
// Once a session is obtained, SessionStart checks for guards.
// If they don't pass it skips the handler.
//
// Finally, if all guards pass, Start is invoked with the session state.
func (session *Session[T]) Start(handler SessionHandler[T]) {
	if !session.Exists() {
		session.Save()
	} else {
		session.Load()
	}

	state := &session.state

	for _, guard := range session.guards {
		passed := false
		guard(session.connection, state, func() { passed = true })
		if !passed {
			if nil == state {
				session.Destroy()
				return
			}
			session.Save()
			return
		}
	}

	handler(state)

	if nil == state {
		session.Destroy()
		return
	}

	session.Save()
}
