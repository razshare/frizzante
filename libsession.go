package frizzante

import (
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
	"strconv"
	"time"
)

type SessionBuilder = func(session *Session)
type Session struct {
	request   *Request
	response  *Response
	get       func(key string) []byte
	set       func(key string, value []byte)
	has       func(key string) bool
	destroy   func()
	onDestroy func()
	id        string
}

func sessionCreate(request *Request, response *Response, builder SessionBuilder) *Session {
	uuidV4, sessionIdError := uuid.NewV4()

	if sessionIdError != nil {
		NotifierSendError(request.server.notifier, sessionIdError)
		return nil
	}

	session := &Session{
		request:  request,
		response: response,
		id:       uuidV4.String(),
	}

	builder(session)

	if nil == session.get {
		session.get = func(key string) []byte { return nil }
	}

	if nil == session.set {
		session.set = func(key string, value []byte) {}
	}

	if nil == session.has {
		session.has = func(key string) bool { return false }
	}

	if nil == session.destroy {
		session.destroy = func() {}
	}

	session.onDestroy = func() {
		delete(sessions, session.id)
	}

	sessions[session.id] = session

	ResponseSendCookie(response, "session-id", session.id)
	return session
}

var sessions = map[string]*Session{}

// SessionStart starts the session and returns its state.
func SessionStart(request *Request, response *Response, builder SessionBuilder) *Session {
	var sessionIdCookie *http.Cookie
	sessionIdCookies := request.httpRequest.CookiesNamed("session-id")
	sessionIdCookiesLen := 0

	for _, cookie := range sessionIdCookies {
		sessionIdCookie = cookie
		sessionIdCookiesLen++
	}

	if 0 == sessionIdCookiesLen || nil == sessionIdCookie {
		// Create new session.
		return sessionCreate(request, response, builder)
	}

	// Try to retrieve session.
	session, sessionExists := sessions[sessionIdCookie.Value]

	if !sessionExists {
		session = &Session{
			request:  request,
			response: response,
			id:       sessionIdCookie.Value,
		}

		builder(session)

		if nil == session.get {
			session.get = func(key string) []byte { return nil }
		}

		if nil == session.set {
			session.set = func(key string, value []byte) {}
		}

		if nil == session.has {
			session.has = func(key string) bool { return false }
		}

		if nil == session.destroy {
			session.destroy = func() {}
		}

		session.onDestroy = func() {
			delete(sessions, session.id)
		}

		sessions[session.id] = session
	}

	return session
}

// SessionId gets the session id.
func SessionId(self *Session) string {
	return self.id
}

// SessionGet gets a property.
func SessionGet(self *Session, key string) []byte {
	return self.get(key)
}

// SessionSet creates or updates a property.
func SessionSet(self *Session, key string, value []byte) {
	self.set(key, value)
}

// SessionGetString gets a property as string.
func SessionGetString(self *Session, key string) string {
	return string(self.get(key))
}

// SessionGetInt64 gets a property as int64.
func SessionGetInt64(self *Session, key string) int64 {
	value, conversionError := strconv.ParseInt(string(self.get(key)), 10, 64)
	if nil != conversionError {
		NotifierSendError(self.request.server.notifier, conversionError)
	}
	return value
}

// SessionSetInt64 creates or updates a property as int64.
func SessionSetInt64(self *Session, key string, value int64) {
	self.set(key, []byte(strconv.FormatInt(value, 10)))
}

// SessionGetFloat64 gets a property as string.
func SessionGetFloat64(self *Session, key string) float64 {
	value, conversionError := strconv.ParseFloat(string(self.get(key)), 64)
	if nil != conversionError {
		NotifierSendError(self.request.server.notifier, conversionError)
	}
	return value
}

// SessionSetFloat64 creates or updates a property as string.
func SessionSetFloat64(self *Session, key string, value float64) {
	self.set(key, []byte(strconv.FormatFloat(value, 'f', -1, 64)))
}

// SessionGetTime gets a property as time.
func SessionGetTime(self *Session, key string) time.Time {
	parsedTime, parseError := time.Parse(time.RFC3339, string(self.get(key)))
	if nil != parseError {
		NotifierSendError(self.request.server.notifier, parseError)
	}
	return parsedTime
}

// SessionSetTime creates or updates a property as time.
func SessionSetTime(self *Session, key string, value time.Time) {
	self.set(key, []byte(value.Format(time.RFC3339)))
}

// SessionGetTimeWithLayout gets a property as time.
func SessionGetTimeWithLayout(self *Session, key string, layout string) time.Time {
	parsedTime, parseError := time.Parse(layout, string(self.get(key)))
	if nil != parseError {
		NotifierSendError(self.request.server.notifier, parseError)
	}
	return parsedTime
}

// SessionSetTimeLayout creates or updates a property as time.
func SessionSetTimeLayout(self *Session, key string, value time.Time, layout string) {
	self.set(key, []byte(value.Format(layout)))
}

// SessionGetJson gets a property as json.
func SessionGetJson[T any](self *Session, key string) T {
	var value T
	data := self.get(key)
	unmarshalError := json.Unmarshal(data, &value)
	if nil != unmarshalError {
		NotifierSendError(self.request.server.notifier, unmarshalError)
	}
	return value
}

// SessionSetString creates or updates a property as string.
func SessionSetString(self *Session, key string, value string) {
	self.set(key, []byte(value))
}

// SessionSetJson creates or updates a property as json.
func SessionSetJson(self *Session, key string, value any) {
	data, marshalError := json.Marshal(value)
	if nil != marshalError {
		NotifierSendError(self.request.server.notifier, marshalError)
		return
	}
	self.set(key, data)
}

// SessionHas checks if the property exists.
func SessionHas(self *Session, key string) bool {
	return self.has(key)
}

// SessionDestroy destroys the session.
func SessionDestroy(self *Session) {
	self.destroy()
}

// SessionWithGetHandler sets the get handler.
func SessionWithGetHandler(self *Session, handler func(key string) []byte) {
	self.get = handler
}

// SessionWithSetHandler sets the set handler.
func SessionWithSetHandler(self *Session, handler func(key string, value []byte)) {
	self.set = handler
}

// SessionWithHasHandler sets the has handler.
func SessionWithHasHandler(self *Session, handler func(key string) bool) {
	self.has = handler
}

// SessionWithDestroyHandler sets the destroy handler.
func SessionWithDestroyHandler(self *Session, handler func()) {
	self.destroy = func() {
		self.onDestroy()
		handler()
	}
}
