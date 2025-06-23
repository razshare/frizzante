package operators

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/globals"
)

// Id tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func (operator *Operator) Id() string {
	if "" != operator.Connection.SessionId {
		return operator.Connection.SessionId
	}

	var sessionId string
	cookies := operator.Connection.Request.CookiesNamed("session-id")
	count := 0

	for _, cookie := range cookies {
		sessionId = cookie.Value
		count++
	}

	if count > 0 {
		operator.Connection.SessionId = sessionId
		return sessionId
	}

	// Create new session.
	id4, idError := uuid.NewV4()
	if idError != nil {
		operator.Connection.Notifier.SendErrorAndTrace(idError, 1)
		operator.Connection.SessionId = ""
		return ""
	}
	sessionId = id4.String()
	operator.Connection.SendCookie("session-id", sessionId)
	operator.Connection.SessionId = sessionId
	return sessionId
}

// Exists checks if the session exists into the archive.
func (operator *Operator) Exists() bool {
	id := operator.Id()
	has, hasError := operator.Archive.Has(id, globals.SessionKey)
	if hasError != nil {
		operator.Connection.Notifier.SendErrorAndTrace(hasError, 1)
	}
	return has
}

// Save saves the session into the archive.
func (operator *Operator) Save(val any) {
	id := operator.Id()
	readBytes, marshalError := json.Marshal(val)
	if marshalError != nil {
		operator.Connection.Notifier.SendErrorAndTrace(marshalError, 1)
		return
	}

	setError := operator.Archive.Set(id, globals.SessionKey, readBytes)
	if setError != nil {
		operator.Connection.Notifier.SendErrorAndTrace(setError, 1)
	}
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func (operator *Operator) Load(val any) {
	id := operator.Id()
	has, hasError := operator.Archive.Has(id, globals.SessionKey)
	if hasError != nil {
		operator.Connection.Notifier.SendErrorAndTrace(hasError, 1)
		return
	}
	if has {
		readBytes, getError := operator.Archive.Get(id, globals.SessionKey)
		if getError != nil {
			operator.Connection.Notifier.SendErrorAndTrace(getError, 1)
			return
		}
		unmarshalError := json.Unmarshal(readBytes, val)
		if unmarshalError != nil {
			operator.Connection.Notifier.SendErrorAndTrace(unmarshalError, 1)
		}
	}
	return
}

// Destroy removes the session from the archive.
func (operator *Operator) Destroy() {
	id := operator.Id()
	destroyError := operator.Archive.RemoveDomain(id)
	if destroyError != nil {
		operator.Connection.Notifier.SendErrorAndTrace(destroyError, 1)
	}
}
