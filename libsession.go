package frizzante

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
	"time"
)

var sessions = map[string]*Session{}

type SessionBuilder = func(session *Session)

type memorySessionStore struct {
	data           map[string]string
	lastActivityAt time.Time
}

type Session struct {
	id       string
	get      func(key string) string
	set      func(key string, value string)
	remove   func(key string)
	has      func(key string) bool
	validate func() (isValid bool)
	destroy  func()
	notifier *Notifier
}

// SessionStart first tries to retrieve the client session, then,
// if the client is not associated with any existing session,
// it will automatically create a new empty session for said client.
//
// It always returns three functions, get, set and unset.
//
// Use get to retrieve a property from the session.
//
// Use set to create a new property or update an existing one to the session.
//
// Use unset to remove a property from the session.
func SessionStart(request *Request, response *Response) *Session {
	var sessionIdCookie *http.Cookie
	sessionIdCookies := request.httpRequest.CookiesNamed("session-id")
	sessionIdCookiesLen := len(sessionIdCookies)

	if 0 == sessionIdCookiesLen {
		uuidV4, sessionIdError := uuid.NewV4()

		if sessionIdError != nil {
			NotifierSendError(request.server.notifier, sessionIdError)
		}

		sessionId := uuidV4.String()
		freshSession := &Session{id: sessionId}
		request.server.sessionBuilder(freshSession)
		sessions[freshSession.id] = freshSession
		ResponseSendCookie(response, "session-id", freshSession.id)
		return freshSession
	}

	var sessionExistsInMemory bool
	var providedSessionId string
	var session *Session

	for _, cookie := range sessionIdCookies {
		providedSessionId = cookie.Value
		session, sessionExistsInMemory = sessions[providedSessionId]
		if sessionExistsInMemory {
			sessionIdCookie = cookie
			break
		}
	}

	if !sessionExistsInMemory {
		var sessionId string
		if "" != providedSessionId {
			sessionId = providedSessionId
		} else {
			uuidV4, sessionIdError := uuid.NewV4()
			if sessionIdError != nil {
				NotifierSendError(request.server.notifier, sessionIdError)
			}

			sessionId = uuidV4.String()
		}

		freshSession := &Session{id: sessionId}
		request.server.sessionBuilder(freshSession)
		sessions[freshSession.id] = freshSession
		ResponseSendCookie(response, "session-id", sessionId)
		return freshSession
	}

	if !session.validate() {
		delete(sessions, sessionIdCookie.Value)
		session.destroy()
		return SessionStart(request, response)
	}

	ResponseSendCookie(response, "session-id", session.id)
	return session
}

// SessionWithId sets the session id.
func SessionWithId(self *Session, id string) {
	self.id = id
}

// SessionId gets the session id.
func SessionId(self *Session) string {
	return self.id
}

// SessionGet gets a property from the session store and unmarshals it as T.
func SessionGet[T any](self *Session, key string) T {
	var value T
	content := self.get(key)

	marshalError := json.Unmarshal([]byte(content), &value)
	if marshalError != nil {
		NotifierSendError(self.notifier, marshalError)
	}

	return value
}

// SessionSet marshals a value and saves it into the session store.
func SessionSet[T any](self *Session, key string, value T) {
	content, marshalError := json.Marshal(value)
	if marshalError != nil {
		NotifierSendError(self.notifier, marshalError)
		return
	}
	self.set(key, string(content))
}

// SessionRemove unsets a property in the session store.
func SessionRemove(self *Session, key string) {
	self.remove(key)
}

// SessionHas checks if a key exists in the session store.
func SessionHas(self *Session, key string) bool {
	return self.has(key)
}

// SessionValidate validates the session.
func SessionValidate(self *Session) bool {
	return self.validate()
}

// SessionDestroy destroys the session.
func SessionDestroy(self *Session) {
	self.destroy()
}

// SessionWithGetter sets the getter, which retrieves a property from the session store.
func SessionWithGetter(self *Session, get func(key string) (value string)) {
	self.get = get
}

// SessionWithSetter sets the setter, which creates or modify a property to the session store.
func SessionWithSetter(self *Session, set func(key string, value string)) {
	self.set = set
}

// SessionWithRemover sets the unsetter, which removes a property from the session store.
func SessionWithRemover(self *Session, remove func(key string)) {
	self.remove = remove
}

// SessionWithChecker sets the key checker, which checks if a key exists in the session store.
func SessionWithChecker(self *Session, check func(key string) bool) {
	self.has = check
}

// SessionWithValidator sets the validator, which validates if the session is still alive.
func SessionWithValidator(self *Session, validator func() (valid bool)) {
	self.validate = validator
}

// SessionWithDestroyer sets the destroyer, which destroys the session.
func SessionWithDestroyer(self *Session, destroyer func()) {
	self.destroy = destroyer
}

// SessionWithNotifier sets the notifier.
func SessionWithNotifier(self *Session, notifier *Notifier) {
	self.notifier = notifier
}

// SessionNotifier gets the notifier.
func SessionNotifier(self *Session) *Notifier {
	return self.notifier
}

// SessionBuilderCreateWithMemory creates a session builder that uses memory (RAM) as a backend.
func SessionBuilderCreateWithMemory() SessionBuilder {
	allData := map[string]map[string]string{}
	operatingGlobal := map[string]map[string]chan int{}
	lastActivities := map[string]time.Time{}
	return func(session *Session) {
		var operating map[string]chan int
		sessionId := SessionId(session)
		data, dataExists := allData[sessionId]
		if !dataExists {
			data = map[string]string{}
			allData[sessionId] = data
		}

		operatingLocal, operatingExists := operatingGlobal[sessionId]
		if operatingExists {
			operating = operatingLocal
		} else {
			operating = map[string]chan int{}
			operatingGlobal[sessionId] = operating
		}

		SessionWithGetter(session, func(key string) string {
			lastActivities[sessionId] = time.Now()
			value := data[key]
			return value
		})

		SessionWithSetter(session, func(key string, value string) {
			lastActivities[sessionId] = time.Now()
			data[key] = value
		})

		SessionWithRemover(session, func(key string) {
			lastActivities[sessionId] = time.Now()
			delete(data, key)
		})

		SessionWithChecker(session, func(key string) bool {
			_, ok := data[key]
			return ok
		})

		SessionWithValidator(session, func() bool {
			lastActivity, ok := lastActivities[sessionId]
			if !ok {
				lastActivity = time.Now()
			}
			return time.Since(lastActivity).Minutes() < 30
		})

		SessionWithDestroyer(session, func() {
			delete(allData, sessionId)
		})
	}
}

// SessionBuilderCreateWithArchive creates a session builder that uses an archive a backend.
func SessionBuilderCreateWithArchive(archive *Archive) SessionBuilder {
	operatingGlobal := map[string]map[string]chan int{}
	lastActivities := map[string]time.Time{}
	return func(session *Session) {
		var operating map[string]chan int
		sessionId := SessionId(session)
		archiveDomain := sessionId

		operatingLocal, operatingExists := operatingGlobal[sessionId]
		if operatingExists {
			operating = operatingLocal
		} else {
			operating = map[string]chan int{}
			operatingGlobal[sessionId] = operating
		}

		lock := func(key string) chan int {
			op, ok := operating[key]
			if !ok {
				op = make(chan int, 1)
				op <- 0
				operating[key] = op
			}

			return op
		}

		SessionWithGetter(session, func(key string) string {
			<-lock(key)
			lastActivities[sessionId] = time.Now()
			value := ArchiveGet(archive, archiveDomain, key)
			lock(key) <- 0
			return value
		})

		SessionWithSetter(session, func(key string, value string) {
			<-lock(key)
			lastActivities[sessionId] = time.Now()
			ArchiveSet(archive, archiveDomain, key, value)
			lock(key) <- 0
		})

		SessionWithRemover(session, func(key string) {
			<-lock(key)
			lastActivities[sessionId] = time.Now()
			ArchiveRemove(archive, archiveDomain, key)
			lock(key) <- 0
		})

		SessionWithChecker(session, func(key string) bool {
			<-lock(key)
			has := ArchiveHas(archive, archiveDomain, key)
			lock(key) <- 0
			return has
		})

		SessionWithValidator(session, func() bool {
			lastActivity, ok := lastActivities[sessionId]
			if !ok {
				lastActivity = time.Now()
			}
			return time.Since(lastActivity).Minutes() < 30
		})

		SessionWithDestroyer(session, func() {
			ArchiveRemoveDomain(archive, archiveDomain)
		})
	}
}
