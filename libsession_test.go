package frizzante

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestSessionStart(test *testing.T) {
	server := ServerCreate()
	port := NextNumber(8080)
	ServerWithPort(server, port)
	ServerWithApiBuilder(server, func(context ApiContext) {
		pattern, handler := context()
		pattern("GET /")
		handler(func(request *Request, response *Response) {
			get, _, _ := SessionStart(request, response)
			name := get("name", "world").(string)
			SendEcho(response, fmt.Sprintf("hello %s", name))
		})
	})
	ServerWithApiBuilder(server, func(context ApiContext) {
		pattern, handler := context()
		pattern("POST /")
		handler(func(request *Request, response *Response) {
			_, set, _ := SessionStart(request, response)
			name := ReceiveMessage(request)
			set("name", name)
			SendEcho(response, "")
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
