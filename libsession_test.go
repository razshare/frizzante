package frizzante

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

type State struct {
	Name string
}

var memory = map[string]State{}

// Memory builds sessions in memory.
func Memory(session *Session[State]) {
	session.WithExistsHandler(func() bool {
		_, exists := memory[session.Id]
		return exists
	})

	session.WithLoadHandler(func() {
		session.Data = memory[session.Id]
	})

	session.WithSaveHandler(func() {
		memory[session.Id] = session.Data
	})

	session.WithDestroyHandler(func() {
		delete(memory, session.Id)
	})

	if session.Exists() {
		session.Load()
		return
	}

	session.Data = State{Name: "world"}
}

func TestSessionStart(test *testing.T) {
	server := NewServer()
	port := NextNumber(8080)
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.OnRequest("GET /", []Guard{}, func(request *Request, response *Response) {
		session := SessionStart(request, response, Memory)
		response.SendMessage(fmt.Sprintf("hello %s", session.Data.Name))
	})

	server.OnRequest("POST /", []Guard{}, func(request *Request, response *Response) {
		session := SessionStart(request, response, Memory)
		session.Data.Name = request.ReceiveMessage()
		response.SendMessage("")
	})

	go server.Start()
	defer server.Stop()

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
