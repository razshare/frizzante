package sessions

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/stack"
)

// New creates a new session with a given initial state.
func New[T any](connection *connections.Connection, initialState T) *Session[T] {
	return &Session[T]{
		Connection: connection,
		State:      &initialState,
	}
}

// Start starts the user's session.
//
// If the user provides a valid "session-id" cookie,
// Start will retrieve the relative session.
//
// If the user doesn't provide a valid "session-id" cookie,
// Start will create a new session along with a new "session-id",
// which it sends to the user as a cookie.
func (session *Session[T]) Start() *Session[T] {
	if !session.Exists() {
		session.Save()
		return session
	}

	session.Load()
	return session
}

// Id tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func (session *Session[T]) Id() string {
	if "" != session.Connection.SessionId {
		return session.Connection.SessionId
	}

	var id string
	cookies := session.Connection.Request.CookiesNamed("session-id")
	connection := 0

	for _, cookie := range cookies {
		id = cookie.Value
		connection++
	}

	if connection > 0 {
		session.Connection.SessionId = id
		return id
	}

	// Create new session.
	idObject, idObjectError := uuid.NewV4()
	if idObjectError != nil {
		session.Connection.SessionId = ""
		session.Connection.ErrorLog.Println(idObjectError, stack.Trace())
		return ""
	}

	id = idObject.String()

	session.Connection.SendCookie("session-id", id)

	session.Connection.SessionId = id

	return id
}

// Exists checks if the session exists into the archive.
func (session *Session[T]) Exists() bool {
	id := session.Id()

	exists, existsError := session.Connection.SessionArchive.Has(id, globals.SessionKey)
	if existsError != nil {
		session.Connection.ErrorLog.Println(existsError, stack.Trace())
		return false
	}
	return exists
}

// Save saves the session into the archive.
func (session *Session[T]) Save() {
	id := session.Id()

	data, jsonError := json.Marshal(session.State)
	if jsonError != nil {
		session.Connection.ErrorLog.Println(jsonError, stack.Trace())
		return
	}

	archiveError := session.Connection.SessionArchive.Set(id, globals.SessionKey, data)
	if archiveError != nil {
		session.Connection.ErrorLog.Println(archiveError, stack.Trace())
	}
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func (session *Session[T]) Load() {
	id := session.Id()

	exists, existsError := session.Connection.SessionArchive.Has(id, globals.SessionKey)
	if existsError != nil {
		session.Connection.ErrorLog.Println(existsError, stack.Trace())
		return
	}

	if exists {
		data, getError := session.Connection.SessionArchive.Get(id, globals.SessionKey)
		if getError != nil {
			session.Connection.ErrorLog.Println(getError, stack.Trace())
			return
		}

		jsonError := json.Unmarshal(data, session.State)
		if jsonError != nil {
			session.Connection.ErrorLog.Println(jsonError, stack.Trace())
			return
		}
	}
}

// Destroy removes the session from the archive.
func (session *Session[T]) Destroy() {
	id := session.Id()

	archiveError := session.Connection.SessionArchive.RemoveDomain(id)
	if archiveError != nil {
		session.Connection.ErrorLog.Println(archiveError, stack.Trace())
		return
	}
}
