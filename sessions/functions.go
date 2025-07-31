package sessions

import (
	"encoding/json"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/actions"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/locks"
	"github.com/razshare/frizzante/traces"
	"os"
	"path/filepath"
)

// New creates a new session with a zero initial state.
func New[T any](connection *connections.Connection, state T) *Session[T] {
	name := filepath.Join(".gen", "sessions")
	lock := locks.New()
	return &Session[T]{
		Get: func(domain string, key string) ([]byte, error) {
			if "" == name {
				return nil, errors.New("disk archive name is blank")
			}

			mutex := locks.FindAndAcquire(lock, domain, key)
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

			mutex := locks.FindAndAcquire(lock, domain, key)
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

			mutex := locks.FindAndAcquire(lock, domain, key)
			mutex.Lock()
			defer mutex.Unlock()

			fileName := filepath.Join(name, domain, key)
			return files.IsFile(fileName), nil
		},
		Remove: func(domain string, key string) error {
			if "" == name {
				return errors.New("disk archive name is blank")
			}

			mutex := locks.FindAndAcquire(lock, domain, key)
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

			mutex := locks.FindAndAcquire(lock, domain)
			mutex.Lock()
			defer mutex.Unlock()

			return files.IsDirectory(filepath.Join(name, domain)), nil
		},
		RemoveDomain: func(domain string) error {
			if "" == name {
				return errors.New("disk archive name is blank")
			}

			mutex := locks.FindAndAcquire(lock, domain)
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
		State:      &state,
	}
}

// Start loads the State if the current connection defines a session-id cookie.
//
// If the session-id cookie is missing it will create a new one and send it to the user.
func Start[T any](self *Session[T]) *Session[T] {

	if !Exists(self) {
		Save(self)
		return self
	}

	Load(self)
	return self
}

// Id tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func Id[T any](self *Session[T]) string {
	if "" != self.Connection.SessionId {
		return self.Connection.SessionId
	}

	var id string
	cookies := self.Connection.Request.CookiesNamed("session-id")
	connection := 0

	for _, cookie := range cookies {
		id = cookie.Value
		connection++
	}

	if connection > 0 {
		self.Connection.SessionId = id
		return id
	}

	// Create new session.
	idObject, idObjectError := uuid.NewV4()
	if idObjectError != nil {
		self.Connection.SessionId = ""
		traces.Trace(self.Connection.Http.ErrorLog, idObjectError)
		return ""
	}

	id = idObject.String()

	actions.SendCookie(self.Connection, "session-id", id)

	self.Connection.SessionId = id

	return id
}

// Exists checks if the session exists into the archive.
func Exists[T any](self *Session[T]) bool {
	id := Id(self)

	exists, existsError := self.Has(id, globals.SessionKey)
	if existsError != nil {
		traces.Trace(self.Connection.Http.ErrorLog, existsError)
		return false
	}
	return exists
}

// Save saves the session into the archive.
func Save[T any](self *Session[T]) {
	id := Id(self)

	data, jsonError := json.Marshal(self.State)
	if jsonError != nil {
		traces.Trace(self.Connection.Http.ErrorLog, jsonError)
		return
	}

	archiveError := self.Set(id, globals.SessionKey, data)
	if archiveError != nil {
		traces.Trace(self.Connection.Http.ErrorLog, archiveError)
	}
}

// Load loads the session from the archive.
//
// If the session is not found in the archive it creates it.
func Load[T any](self *Session[T]) {
	id := Id(self)

	exists, existsError := self.Has(id, globals.SessionKey)
	if existsError != nil {
		traces.Trace(self.Connection.Http.ErrorLog, existsError)
		return
	}

	if exists {
		data, getError := self.Get(id, globals.SessionKey)
		if getError != nil {
			traces.Trace(self.Connection.Http.ErrorLog, getError)
			return
		}

		jsonError := json.Unmarshal(data, self.State)
		if jsonError != nil {
			traces.Trace(self.Connection.Http.ErrorLog, jsonError)
			return
		}
	}
}

// Destroy removes the session from the archive.
func Destroy[T any](self *Session[T]) {
	id := Id(self)

	archiveError := self.RemoveDomain(id)
	if archiveError != nil {
		traces.Trace(self.Connection.Http.ErrorLog, archiveError)
		return
	}
}
