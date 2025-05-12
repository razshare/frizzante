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
	load     func()
	save     func()
	validate func() bool
	destroy  func()
	Id       string
	State    T
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

	if nil == session.load {
		session.load = func() {}
	}

	if nil == session.save {
		session.save = func() {}
	}

	if nil == session.validate {
		session.validate = func() bool { return true }
	}

	if nil == session.destroy {
		session.destroy = func() {}
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
		return sessionCreate[T](request, response, builder).State
	}

	// Retrieve session.
	session := &Session[T]{
		request:  request,
		response: response,
		Id:       sessionIdCookie.Value,
	}

	builder(session)

	if nil == session.load {
		session.load = func() {}
	}

	if nil == session.save {
		session.save = func() {}
	}

	if nil == session.validate {
		session.validate = func() bool { return true }
	}

	if nil == session.destroy {
		session.destroy = func() {}
	}

	session.load()
	if session.validate() {
		response.after = append(response.after, func() {
			session.save()
		})
		return session.State
	}

	session.destroy()

	return sessionCreate[T](request, response, builder).State
}

// SessionWithLoadHandler sets the load handler.
func SessionWithLoadHandler[T any](self *Session[T], loader func()) {
	self.load = loader
}

// SessionWithSaveHandler sets the save handler.
func SessionWithSaveHandler[T any](self *Session[T], handler func()) {
	self.save = handler
}

// SessionWithValidateHandler sets the validate handler.
func SessionWithValidateHandler[T any](self *Session[T], handler func() bool) {
	self.validate = handler
}

// SessionWithDestroyHandler sets the destroy handler.
func SessionWithDestroyHandler[T any](self *Session[T], handler func()) {
	self.destroy = handler
}
