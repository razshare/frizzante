package sessions

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/servers"
	"github.com/razshare/frizzante/traces"
	"path/filepath"
)

// New creates a new session with a zero initial state.
func New[T any](connection *servers.Connection, state T) *Session[T] {
	return &Session[T]{
		Archive:    archives.NewDiskArchive(filepath.Join(".gen", "sessions")),
		Connection: connection,
		State:      &state,
	}
}

// Start loads the State if the current connection defines a session-id cookie.
//
// If the session-id cookie is missing it will create a new one and send it to the user.
func (session *Session[T]) Start() {
	exists := session.Exists()

	if !exists {
		session.Save()
		return
	}

	session.Load()
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
		traces.Trace(session.Connection.Web.ErrorLog, idObjectError)
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

	exists, existsError := session.Archive.Has(id, globals.SessionKey)
	if existsError != nil {
		traces.Trace(session.Connection.Web.ErrorLog, existsError)
		return false
	}
	return exists
}

// Save saves the session into the archive.
func (session *Session[T]) Save() {
	id := session.Id()

	data, jsonError := json.Marshal(session.State)
	if jsonError != nil {
		traces.Trace(session.Connection.Web.ErrorLog, jsonError)
		return
	}

	archiveError := session.Archive.Set(id, globals.SessionKey, data)
	if archiveError != nil {
		traces.Trace(session.Connection.Web.ErrorLog, archiveError)
	}
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func (session *Session[T]) Load() {
	id := session.Id()

	exists, existsError := session.Archive.Has(id, globals.SessionKey)
	if existsError != nil {
		traces.Trace(session.Connection.Web.ErrorLog, existsError)
		return
	}

	if exists {
		data, getError := session.Archive.Get(id, globals.SessionKey)
		if getError != nil {
			traces.Trace(session.Connection.Web.ErrorLog, getError)
			return
		}

		jsonError := json.Unmarshal(data, session.State)
		if jsonError != nil {
			traces.Trace(session.Connection.Web.ErrorLog, jsonError)
			return
		}
	}
}

// Destroy removes the session from the archive.
func (session *Session[T]) Destroy() {
	id := session.Id()

	archiveError := session.Archive.RemoveDomain(id)
	if archiveError != nil {
		traces.Trace(session.Connection.Web.ErrorLog, archiveError)
		return
	}
}
