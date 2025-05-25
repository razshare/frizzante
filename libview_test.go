package frizzante

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type Data struct {
	Name string `json:"name"`
}

func TestRenderServer(test *testing.T) {
	port := NextNumber(8080)
	notifier := NewNotifier()
	server := NewServer()
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.WithNotifier(notifier)
	server.WithRoute("GET /welcome", func(req *Request, res *Response) {
		res.SendView(View{
			Name:       "Welcome",
			RenderMode: RenderModeServer,
			Data:       Data{Name: "world"},
		})
	})

	go server.Start(efs)
	defer server.Stop()
	time.Sleep(1 * time.Second)

	expected := "<h1>Hello world.</h1>"
	actual, getError := HttpGet(fmt.Sprintf("http://127.0.0.1:%d/welcome", port), nil)
	if getError != nil {
		test.Fatal(getError)
	}

	ok := strings.Contains(actual, expected)

	if !ok {
		test.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}

func TestRenderClient(test *testing.T) {
	port := NextNumber(8080)
	notifier := NewNotifier()
	server := NewServer()
	server.WithNotifier(notifier)
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.WithRoute("GET /welcome", func(req *Request, res *Response) {
		res.SendView(View{
			Name:       "Welcome",
			RenderMode: RenderModeClient,
			Data:       Data{Name: "world"},
		})
	})

	go server.Start(efs)
	defer server.Stop()
	time.Sleep(1 * time.Second)

	expected := "<script type=\"application/javascript\">function target(){return document.getElementById("
	actual, getError := HttpGet(fmt.Sprintf("http://127.0.0.1:%d/welcome", port), nil)
	if getError != nil {
		test.Fatal(getError)
	}

	ok := strings.Contains(actual, expected)

	if !ok {
		test.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}
