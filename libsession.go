package frizzante

import (
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
)

type SessionInitializer[T any] = func() T
type SessionBuilder[T any] = func(session *Session[T])
type Session[T any] struct {
	request  *Request
	response *Response
	validate func() bool
	destroy  func()
	load     func()
	save     func()
	Id       string
	Store    T
}

func sessionCreate[T any](request *Request, response *Response, builder SessionBuilder[T]) *Session[T] {
	uuidV4, sessionIdError := uuid.NewV4()

	if sessionIdError != nil {
		NotifierSendError(request.server.notifier, sessionIdError)
		return nil
	}

	sessionId := uuidV4.String()
	session := &Session[T]{
		request:  request,
		response: response,
		Id:       sessionId,
	}

	response.after = append(response.after, func() {
		session.save()
	})

	builder(session)

	if nil == session.save {
		session.save = func() {}
	}

	if nil == session.load {
		session.load = func() {}
	}

	if nil == session.destroy {
		session.destroy = func() {}
	}

	if nil == session.validate {
		session.validate = func() bool { return true }
	}

	session.load()
	ResponseSendCookie(response, "session-id", session.Id)
	return session
}

// SessionStart starts the session and returns its state.
func SessionStart[T any](request *Request, response *Response, builder SessionBuilder[T]) T {
	var sessionIdCookie *http.Cookie
	sessionIdCookies := request.httpRequest.CookiesNamed("session-id")
	sessionIdCookiesLen := 0

	for _, cookie := range sessionIdCookies {
		sessionIdCookie = cookie
		sessionIdCookiesLen++
	}

	if 0 == sessionIdCookiesLen || nil == sessionIdCookie {
		// Create new session.
		return sessionCreate[T](request, response, builder).Store
	}

	// Retrieve session.
	session := &Session[T]{
		request:  request,
		response: response,
		Id:       sessionIdCookie.Value,
	}

	builder(session)

	if nil == session.save {
		session.save = func() {}
	}

	if nil == session.load {
		session.load = func() {}
	}

	if nil == session.destroy {
		session.destroy = func() {}
	}

	if nil == session.validate {
		session.validate = func() bool { return true }
	}

	if nil != session.load {
		session.load()
	}

	if nil != session.validate && session.validate() {
		response.after = append(response.after, func() {
			session.save()
		})
		return session.Store
	}

	if nil != session.destroy {
		session.destroy()
	}

	return sessionCreate[T](request, response, builder).Store
}

// SessionWithLoader sets the loader.
func SessionWithLoader[T any](self *Session[T], loader func()) {
	self.load = loader
}

// SessionWithValidator sets the validator.
func SessionWithValidator[T any](self *Session[T], validator func() bool) {
	self.validate = validator
}

// SessionWithDestroyer sets the destroyer.
func SessionWithDestroyer[T any](self *Session[T], destroyer func()) {
	self.destroy = destroyer
}

// SessionWithSaver sets the saver.
func SessionWithSaver[T any](self *Session[T], saver func()) {
	self.save = saver
}
