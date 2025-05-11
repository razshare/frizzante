package frizzante

import (
	"encoding/json"
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

// SessionStart starts the session and returns its state.
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

// SessionBuilderCreate creates a session builder backed by the disk.
func SessionBuilderCreate[T any](archive *Archive, initialize func() *T) SessionBuilder[T] {
	var key = "session.json"
	var operating = map[string]chan int{}
	var sessions = map[string]*T{}
	return func(session *Session[T]) {
		SessionWithLoader(session, func() {
			_, sessionExists := sessions[session.Id]
			if !sessionExists {
				sessions[session.Id] = initialize()
				operating[session.Id] = make(chan int, 1)
				operating[session.Id] <- 0

				<-operating[session.Id]
				if !ArchiveHas(archive, session.Id, key) {
					readBytes, marshalError := json.Marshal(sessions[session.Id])
					if nil != marshalError {
						NotifierSendError(archive.notifier, marshalError)
						operating[session.Id] <- 0
						return
					}
					ArchiveSet(archive, session.Id, key, readBytes)
				}
				operating[session.Id] <- 0
			}

			<-operating[session.Id]
			session.Value = sessions[session.Id]
			readBytes := ArchiveGet(archive, session.Id, key)
			marshalError := json.Unmarshal(readBytes, sessions[session.Id])
			if nil != marshalError {
				NotifierSendError(archive.notifier, marshalError)
			}
			operating[session.Id] <- 0
		})

		SessionWithValidator(session, func() bool {
			return true
		})

		SessionWithSaver(session, func() {
			<-operating[session.Id]
			readBytes, marshalError := json.Marshal(sessions[session.Id])
			if nil != marshalError {
				NotifierSendError(archive.notifier, marshalError)
				operating[session.Id] <- 0
				return
			}
			ArchiveSet(archive, session.Id, key, readBytes)
			operating[session.Id] <- 0
		})

		SessionWithDestroyer(session, func() {
			<-operating[session.Id]
			delete(operating, session.Id)
			operating[session.Id] <- 0
		})
	}
}
