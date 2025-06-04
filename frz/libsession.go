package frz

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
)

type SessionOperator interface {
	Id() string
	Exists() bool
	Load(state any)
	Save(state any)
	Destroy()
	WithArchive(archive Archive) SessionOperator
}

type SessionHandler[T any] = func(state *T)

type ConnectedSessionOperator struct {
	archive    Archive
	connection *Connection
}

// Session creates a new session.
func Session[T any](c *Connection, state T) (*T, SessionOperator) {
	archive := NewDiskArchive().WithName("sessions")
	manager := &ConnectedSessionOperator{
		archive:    archive,
		connection: c,
	}

	if !manager.Exists() {
		manager.Save(&state)
	} else {
		manager.Load(&state)
	}

	return &state, manager
}

// WithArchive sets the archive.
func (session *ConnectedSessionOperator) WithArchive(archive Archive) SessionOperator {
	session.archive = archive
	return session
}

// Id tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func (session *ConnectedSessionOperator) Id() string {
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
		session.connection.server.notifier.SendErrorAndTrace(idError, 1)
		session.connection.sessionId = ""
		return ""
	}
	sessionId = id4.String()
	session.connection.SendCookie("session-id", sessionId)
	session.connection.sessionId = sessionId
	return sessionId
}

// Exists checks if the session exists into the archive.
func (session *ConnectedSessionOperator) Exists() bool {
	id := session.Id()
	has, hasError := session.archive.Has(id, SessionKey)
	if hasError != nil {
		session.connection.server.notifier.SendErrorAndTrace(hasError, 1)
	}
	return has
}

// Save saves the session into the archive.
func (session *ConnectedSessionOperator) Save(state any) {
	id := session.Id()
	readBytes, marshalError := json.Marshal(state)
	if marshalError != nil {
		session.connection.server.notifier.SendErrorAndTrace(marshalError, 1)
		return
	}
	setError := session.archive.Set(id, SessionKey, readBytes)
	if setError != nil {
		session.connection.server.notifier.SendErrorAndTrace(setError, 1)
	}
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func (session *ConnectedSessionOperator) Load(state any) {
	id := session.Id()

	has, hasError := session.archive.Has(id, SessionKey)
	if hasError != nil {
		session.connection.server.notifier.SendErrorAndTrace(hasError, 1)
		return
	}
	if has {
		readBytes, getError := session.archive.Get(id, SessionKey)
		if getError != nil {
			session.connection.server.notifier.SendErrorAndTrace(getError, 1)
			return
		}
		unmarshalError := json.Unmarshal(readBytes, state)
		if unmarshalError != nil {
			session.connection.server.notifier.SendErrorAndTrace(unmarshalError, 1)
		}
	}
	return
}

// Destroy removes the session from the archive.
func (session *ConnectedSessionOperator) Destroy() {
	id := session.Id()
	destroyError := session.archive.RemoveDomain(id)
	if destroyError != nil {
		session.connection.server.notifier.SendErrorAndTrace(destroyError, 1)
	}
}
