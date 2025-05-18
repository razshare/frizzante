package frizzante

import (
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
)

type SessionBuilder[T any] = func(session *Session[T])
type Session[T any] struct {
	request   *Request
	response  *Response
	exists    func() bool
	load      func()
	save      func()
	destroy   func()
	onDestroy func()
	Id        string
	Data      T
}

var sessions = map[string]any{}

func newSessionInitializedAndBuilt[T any](request *Request, response *Response, builder SessionBuilder[T]) *Session[T] {
	uuidV4, sessionIdError := uuid.NewV4()

	if sessionIdError != nil {
		request.server.notifier.SendError(sessionIdError)
		return nil
	}

	session := &Session[T]{
		request:  request,
		response: response,
		Id:       uuidV4.String(),
	}

	builder(session)

	session.onDestroy = func() {
		delete(sessions, session.Id)
	}

	sessions[session.Id] = session

	response.SendCookie("session-id", session.Id)
	return session
}

// SessionStart starts the session and returns it.
//
// SessionStart reads the request's cookies and tries to find a "session-id" cookie.
// If such a cookie references an existing session, it returns that session.
//
// Otherwise, SessionStart automatically creates a new session and returns it instead.
// The new session does not include any data from the old session.
//
// Finally, SessionStart modifies the response by applying a "session-id" cookie,
// referencing the session that's currently being used by the server.
//
// This means there can be cases where a client sends a "session-id" cookie of value "AAA"
// but the server responds with a cookie "session-id" of value "BBB", meaning the client's
// "AAA" session doesn't exist, thus the client should use session "BBB" instead.
func SessionStart[T any](request *Request, response *Response, builder SessionBuilder[T]) *Session[T] {
	var sessionIdCookie *http.Cookie
	sessionIdCookies := request.httpRequest.CookiesNamed("session-id")
	sessionIdCookiesLen := 0

	for _, cookie := range sessionIdCookies {
		sessionIdCookie = cookie
		sessionIdCookiesLen++
	}

	if 0 == sessionIdCookiesLen || nil == sessionIdCookie {
		// Create new session.
		return newSessionInitializedAndBuilt[T](request, response, builder)
	}

	// Try to retrieve session.
	sessionAny, sessionExists := sessions[sessionIdCookie.Value]

	if sessionExists {
		return sessionAny.(*Session[T])
	}

	session := &Session[T]{
		request:  request,
		response: response,
		Id:       sessionIdCookie.Value,
	}

	builder(session)

	session.onDestroy = func() {
		delete(sessions, session.Id)
	}

	sessions[session.Id] = session

	return session
}

// Exists checks if the session exists.
func (session *Session[T]) Exists() bool {
	return session.exists()
}

// Load loads the session.
func (session *Session[T]) Load() {
	session.load()
}

// Save saves the session.
func (session *Session[T]) Save() {
	session.save()
}

// Destroy destroys the session.
func (session *Session[T]) Destroy() {
	session.destroy()
}

// WithExistsHandler sets the exists handler.
func (session *Session[T]) WithExistsHandler(handler func() bool) {
	session.exists = handler
}

// WithLoadHandler sets the load handler.
func (session *Session[T]) WithLoadHandler(handler func()) {
	session.load = handler
}

// WithSaveHandler sets the save handler.
func (session *Session[T]) WithSaveHandler(handler func()) {
	session.save = handler
}

// WithDestroyHandler sets the destroy handler.
func (session *Session[T]) WithDestroyHandler(handler func()) {
	session.destroy = func() {
		session.onDestroy()
		handler()
	}
}
