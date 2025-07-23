package sessions

import (
	"encoding/json"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/roads"
	"os"
	"path/filepath"
)

var Identity = filepath.Join(".gen/sessions")
var Road = roads.New()

// Set sets a value to session.
var Set = func(id string, key string, value []byte) error {
	if "" == Identity {
		return errors.New("sessions identity is blank")
	}

	lane := Road.WithLane(id, key)
	lane.Lock()
	defer lane.Unlock()
	directoryName := filepath.Join(Identity, id)
	if !files.IsDirectory(directoryName) {
		mkdirError := os.MkdirAll(directoryName, os.ModePerm)
		if nil != mkdirError {
			return mkdirError
		}
	}
	fileName := filepath.Join(directoryName, key)
	writeError := os.WriteFile(fileName, value, os.ModePerm)
	if nil != writeError {
		return writeError
	}
	return nil
}

// Get gets a value from the session.
var Get = func(id string, key string) ([]byte, error) {
	if "" == Identity {
		return nil, errors.New("sessions identity is blank")
	}
	lane := Road.WithLane(id, key)
	lane.Lock()
	defer lane.Unlock()
	fileName := filepath.Join(Identity, id, key)
	value, readError := os.ReadFile(fileName)
	if nil != readError {
		return nil, readError
	}
	return value, nil
}

// Has checks if a sessions has a key.
var Has = func(id string, key string) (bool, error) {
	if "" == Identity {
		return false, errors.New("sessions identity is blank")
	}
	lane := Road.WithLane(id, key)
	lane.Lock()
	defer lane.Unlock()
	fileName := filepath.Join(Identity, id, key)
	return files.IsFile(fileName), nil
}

// Remove removes a value from a session.
var Remove = func(id string, key string) error {
	if "" == Identity {
		return errors.New("sessions identity is blank")
	}
	lane := Road.WithLane(id, key)
	lane.Lock()
	defer lane.Unlock()
	fileName := filepath.Join(Identity, id, key)
	removeError := os.Remove(fileName)
	if nil != removeError {
		return removeError
	}
	return nil
}

// HasId checks if a session exists.
var HasId = func(id string) (bool, error) {
	if "" == Identity {
		return false, errors.New("sessions identity is blank")
	}
	lane := Road.WithLane(id)
	lane.Lock()
	defer lane.Unlock()
	return files.IsDirectory(filepath.Join(Identity, id)), nil
}

// RemoveId removes a session.
var RemoveId = func(id string) error {
	if "" == Identity {
		return errors.New("sessions identity is blank")
	}
	lane := Road.WithLane(id)
	lane.Lock()
	defer lane.Unlock()
	directoryName := filepath.Join(Identity, id)
	removeError := os.RemoveAll(directoryName)
	if nil != removeError {
		return removeError
	}
	return nil
}

// StartWith starts a session with a given initial state.
func StartWith[T any](con *connections.Connection, state T) *Session[T] {
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
func Start[T any](con *connections.Connection) *Session[T] {
	var state T
	return StartWith(con, state)
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
	has, hasError := Has(id, globals.SessionKey)
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

	setError := Set(id, globals.SessionKey, readBytes)
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
	has, hasError := Has(id, globals.SessionKey)
	if hasError != nil {
		session.Connection.Notifier.SendErrorAndTrace(hasError, 1)
		return session
	}
	if has {
		readBytes, getError := Get(id, globals.SessionKey)
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
	destroyError := RemoveId(id)
	if destroyError != nil {
		session.Connection.Notifier.SendErrorAndTrace(destroyError, 1)
	}
	return session
}
