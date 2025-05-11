package frizzante

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

type state struct {
	name string
}

var sessions = map[string]*state{}
var operating = map[string]chan int{}

func memory(session *Session[state]) {
	destroyed := false
	SessionWithLoader(session, func() {
		if destroyed {
			return
		}

		_, sessionExists := sessions[session.Id]
		if !sessionExists {
			sessions[session.Id] = &state{}
			operating[session.Id] = make(chan int, 1)
			operating[session.Id] <- 0
		}

		<-operating[session.Id]
		session.Value = sessions[session.Id]
		operating[session.Id] <- 0
	})

	SessionWithValidator(session, func() bool {
		if destroyed {
			return false
		}

		return true
	})

	SessionWithSaver(session, func() {
		if destroyed {
			return
		}

		<-operating[session.Id]
		sessions[session.Id] = session.Value
		operating[session.Id] <- 0
	})

	SessionWithDestroyer(session, func() {
		if destroyed {
			return
		}

		<-operating[session.Id]
		delete(sessions, session.Id)
		destroyed = true
		operating[session.Id] <- 0
	})
}

func TestSessionStart(test *testing.T) {
	server := ServerCreate()
	port := NextNumber(8080)
	ServerWithPort(server, port)
	ServerWithSessionBuilder[state](server, memory)
	ServerWithApiBuilder(server, func(api *Api) {
		ApiWithPattern(api, "GET /")
		ApiWithRequestHandler(api, func(request *Request, response *Response) {
			session := SessionStart[state](request, response)

			if "" == session.name {
				session.name = "world"
			}

			ResponseSendMessage(response, fmt.Sprintf("hello %s", session.name))
		})
	})
	ServerWithApiBuilder(server, func(api *Api) {
		ApiWithPattern(api, "POST /")
		ApiWithRequestHandler(api, func(request *Request, response *Response) {
			session := SessionStart[state](request, response)
			session.name = RequestReceiveMessage(request)
			ResponseSendMessage(response, "")
		})
	})

	go ServerStart(server)
	defer ServerStop(server)

	time.Sleep(1 * time.Second)

	sessionId := ""
	expected1 := "hello world"
	expected2 := "hello test"

	response1, error1 := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if error1 != nil {
		test.Fatal(error1)
	}

	cookies := response1.Cookies()
	for _, cookie := range cookies {
		if "session-id" == cookie.Name {
			sessionId = cookie.Value
			break
		}
	}

	if "" == sessionId {
		test.Fatal("session id not found")
	}

	actualBytes1, readAllError := io.ReadAll(response1.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}
	actual1 := string(actualBytes1)

	if expected1 != actual1 {
		test.Fatal(fmt.Sprintf("Message was expected to be `%s`, received `%s` instead.", expected1, actual1))
	}

	_, postError := HttpPost(fmt.Sprintf("http://127.0.0.1:%d/", port), "test", map[string]string{
		"Cookie": fmt.Sprintf("session-id=%s", sessionId),
	})
	if postError != nil {
		test.Fatal(postError)
	}

	actual2, getError2 := HttpGet(fmt.Sprintf("http://127.0.0.1:%d/", port), map[string]string{
		"Cookie": fmt.Sprintf("session-id=%s", sessionId),
	})
	if getError2 != nil {
		test.Fatal(getError2)
	}

	if expected2 != actual2 {
		test.Fatal(fmt.Sprintf("Message was expected to be `%s`, received `%s` instead.", expected2, actual2))
	}
}
