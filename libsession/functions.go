package libsession

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/libarchive"
	"github.com/razshare/frizzante/libcon"
	"github.com/razshare/frizzante/libglobals"
)

// Session creates a new session.
func Session[T any](c *libcon.Connection, state T) (*T, SessionOperator) {
	archive := libarchive.NewDiskArchive().WithName(".gen/sessions")
	manager := &ConnectedSessionOperator{
		Archive:    archive,
		Connection: c,
	}

	if !manager.Exists() {
		manager.Save(&state)
	} else {
		manager.Load(&state)
	}

	return &state, manager
}

// WithArchive sets the archive.
func (session *ConnectedSessionOperator) WithArchive(archive libarchive.Archive) SessionOperator {
	session.Archive = archive
	return session
}

// Id tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func (session *ConnectedSessionOperator) Id() string {
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
func (session *ConnectedSessionOperator) Exists() bool {
	id := session.Id()
	has, hasError := session.Archive.Has(id, libglobals.SessionKey)
	if hasError != nil {
		session.Connection.Notifier.SendErrorAndTrace(hasError, 1)
	}
	return has
}

// Save saves the session into the archive.
func (session *ConnectedSessionOperator) Save(state any) {
	id := session.Id()
	readBytes, marshalError := json.Marshal(state)
	if marshalError != nil {
		session.Connection.Notifier.SendErrorAndTrace(marshalError, 1)
		return
	}
	setError := session.Archive.Set(id, libglobals.SessionKey, readBytes)
	if setError != nil {
		session.Connection.Notifier.SendErrorAndTrace(setError, 1)
	}
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func (session *ConnectedSessionOperator) Load(state any) {
	id := session.Id()

	has, hasError := session.Archive.Has(id, libglobals.SessionKey)
	if hasError != nil {
		session.Connection.Notifier.SendErrorAndTrace(hasError, 1)
		return
	}
	if has {
		readBytes, getError := session.Archive.Get(id, libglobals.SessionKey)
		if getError != nil {
			session.Connection.Notifier.SendErrorAndTrace(getError, 1)
			return
		}
		unmarshalError := json.Unmarshal(readBytes, state)
		if unmarshalError != nil {
			session.Connection.Notifier.SendErrorAndTrace(unmarshalError, 1)
		}
	}
	return
}

// Destroy removes the session from the archive.
func (session *ConnectedSessionOperator) Destroy() {
	id := session.Id()
	destroyError := session.Archive.RemoveDomain(id)
	if destroyError != nil {
		session.Connection.Notifier.SendErrorAndTrace(destroyError, 1)
	}
}
