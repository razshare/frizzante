package sessions

import (
	"encoding/json"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/locks"
	"github.com/razshare/frizzante/traces"
	"os"
	"path/filepath"
)

// New creates a new session with a given initial state.
func New[T any](connection *connections.Connection, initialState T) *Session[T] {
	name := filepath.Join(".gen", "sessions")
	lock := locks.New()
	return &Session[T]{
		Get: func(domain string, key string) ([]byte, error) {
			if "" == name {
				return nil, errors.New("disk archive name is blank")
			}

			mutex := lock.Acquire(domain, key)
			mutex.Lock()
			defer mutex.Unlock()

			fileName := filepath.Join(name, domain, key)
			value, readError := os.ReadFile(fileName)
			if readError != nil {
				return nil, readError
			}
			return value, nil
		},
		Set: func(domain string, key string, value []byte) error {
			if "" == name {
				return errors.New("disk archive name is blank")
			}

			mutex := lock.Acquire(domain, key)
			mutex.Lock()
			defer mutex.Unlock()

			directoryName := filepath.Join(name, domain)
			if !files.IsDirectory(directoryName) {
				mkdirError := os.MkdirAll(directoryName, os.ModePerm)
				if mkdirError != nil {
					return mkdirError
				}
			}
			fileName := filepath.Join(directoryName, key)
			writeError := os.WriteFile(fileName, value, os.ModePerm)
			if writeError != nil {
				return writeError
			}
			return nil
		},
		Has: func(domain string, key string) (bool, error) {
			if "" == name {
				return false, errors.New("disk archive name is blank")
			}

			mutex := lock.Acquire(domain, key)
			mutex.Lock()
			defer mutex.Unlock()

			fileName := filepath.Join(name, domain, key)
			return files.IsFile(fileName), nil
		},
		Remove: func(domain string, key string) error {
			if "" == name {
				return errors.New("disk archive name is blank")
			}

			mutex := lock.Acquire(domain, key)
			mutex.Lock()
			defer mutex.Unlock()

			fileName := filepath.Join(name, domain, key)
			removeError := os.Remove(fileName)
			if removeError != nil {
				return removeError
			}
			return nil
		},
		HasDomain: func(domain string) (bool, error) {
			if "" == name {
				return false, errors.New("disk archive name is blank")
			}

			mutex := lock.Acquire(domain)
			mutex.Lock()
			defer mutex.Unlock()

			return files.IsDirectory(filepath.Join(name, domain)), nil
		},
		RemoveDomain: func(domain string) error {
			if "" == name {
				return errors.New("disk archive name is blank")
			}

			mutex := lock.Acquire(domain)
			mutex.Lock()
			defer mutex.Unlock()

			directoryName := filepath.Join(name, domain)
			removeError := os.RemoveAll(directoryName)
			if removeError != nil {
				return removeError
			}
			return nil
		},
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
		traces.Trace(session.Connection.ErrorLog, idObjectError)
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

	exists, existsError := session.Has(id, globals.SessionKey)
	if existsError != nil {
		traces.Trace(session.Connection.ErrorLog, existsError)
		return false
	}
	return exists
}

// Save saves the session into the archive.
func (session *Session[T]) Save() {
	id := session.Id()

	data, jsonError := json.Marshal(session.State)
	if jsonError != nil {
		traces.Trace(session.Connection.ErrorLog, jsonError)
		return
	}

	archiveError := session.Set(id, globals.SessionKey, data)
	if archiveError != nil {
		traces.Trace(session.Connection.ErrorLog, archiveError)
	}
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func (session *Session[T]) Load() {
	id := session.Id()

	exists, existsError := session.Has(id, globals.SessionKey)
	if existsError != nil {
		traces.Trace(session.Connection.ErrorLog, existsError)
		return
	}

	if exists {
		data, getError := session.Get(id, globals.SessionKey)
		if getError != nil {
			traces.Trace(session.Connection.ErrorLog, getError)
			return
		}

		jsonError := json.Unmarshal(data, session.State)
		if jsonError != nil {
			traces.Trace(session.Connection.ErrorLog, jsonError)
			return
		}
	}
}

// Destroy removes the session from the archive.
func (session *Session[T]) Destroy() {
	id := session.Id()

	archiveError := session.RemoveDomain(id)
	if archiveError != nil {
		traces.Trace(session.Connection.ErrorLog, archiveError)
		return
	}
}
