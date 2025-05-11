package frizzante

import (
	uuid "github.com/nu7hatch/gouuid"
	"net/http"
)

type SessionBuilder[T any] = func(session *Session[T])
type Session[T any] struct {
	request  *Request
	response *Response
	validate func() bool
	destroy  func()
	load     func()
	save     func()
	resume   func() (*T, bool)
	Id       string
	Value    *T
}

func sessionCreate[T any](request *Request, response *Response) *Session[T] {
	uuidV4, sessionIdError := uuid.NewV4()

	if sessionIdError != nil {
		NotifierSendError(request.server.notifier, sessionIdError)
		return nil
	}

	sessionId := uuidV4.String()
	session := &Session[T]{
		request:  request,
		response: response,
		destroy:  func() {},
		validate: func() bool { return true },
		load:     func() {},
		save:     func() {},
		Id:       sessionId,
	}

	response.after = append(response.after, func() {
		session.save()
	})

	if nil != response.server.sessionBuilder {
		builder := response.server.sessionBuilder.(SessionBuilder[T])
		builder(session)
	}

	session.load()
	ResponseSendCookie(response, "session-id", session.Id)
	return session
}

func SessionStart[T any](request *Request, response *Response) *T {
	var sessionIdCookie *http.Cookie
	sessionIdCookies := request.httpRequest.CookiesNamed("session-id")
	sessionIdCookiesLen := 0

	for _, cookie := range sessionIdCookies {
		sessionIdCookie = cookie
		sessionIdCookiesLen++
	}

	if 0 == sessionIdCookiesLen || nil == sessionIdCookie {
		// Create new session.
		return sessionCreate[T](request, response).Value
	}

	// Retrieve session.
	session := &Session[T]{
		request:  request,
		response: response,
		destroy:  func() {},
		validate: func() bool { return true },
		load:     func() {},
		save:     func() {},
		Id:       sessionIdCookie.Value,
	}

	if nil != response.server.sessionBuilder {
		builder := response.server.sessionBuilder.(SessionBuilder[T])
		builder(session)
	}

	session.load()
	if session.validate() {
		response.after = append(response.after, func() {
			session.save()
		})
		return session.Value
	}

	// Session doesn't exist or is not valid, create a new one.
	session.destroy()
	return sessionCreate[T](request, response).Value
}

func SessionWithLoader[T any](self *Session[T], loader func()) {
	self.load = loader
}

func SessionWithValidator[T any](self *Session[T], validator func() bool) {
	self.validate = validator
}

func SessionWithDestroyer[T any](self *Session[T], destroyer func()) {
	self.destroy = destroyer
}

func SessionWithSaver[T any](self *Session[T], saver func()) {
	self.save = saver
}
