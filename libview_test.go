package frizzante

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type WelcomeData struct {
	Name string `json:"name"`
}

// Ssr.
type SsrController struct {
	PageController
}

func (_ SsrController) Configure() PageConfiguration {
	return PageConfiguration{
		Path: "/",
	}
}

func (_ SsrController) Base(request *Request, response *Response) {
	view := NewView(WelcomeData{Name: "world"})
	view.RenderMode = RenderModeServer
	response.SendView(view)
}

func (_ SsrController) Action(request *Request, response *Response) {
	view := NewView(WelcomeData{Name: "world"})
	view.RenderMode = RenderModeServer
	response.SendView(view)
}

// Csr.
type CsrController struct {
	PageController
}

func (_ CsrController) Configure() PageConfiguration {
	return PageConfiguration{
		Path: "/",
	}
}

func (_ CsrController) Base(request *Request, response *Response) {
	view := NewView(WelcomeData{Name: "world"})
	view.RenderMode = RenderModeClient
	response.SendView(view)
}

func (_ CsrController) Action(request *Request, response *Response) {
	view := NewView(WelcomeData{Name: "world"})
	view.RenderMode = RenderModeClient
	response.SendView(view)
}

func TestRenderServer(test *testing.T) {
	server := NewServer()
	notifier := NewNotifier()
	port := NextNumber(8080)
	server.WithPort(port)
	server.WithHostName("127.0.0.1")
	server.WithNotifier(notifier)
	server.WithEmbeddedFileSystem(&embeddedFileSystem)
	server.WithPageController(SsrController{})

	go server.Start()
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
	server := NewServer()
	notifier := NewNotifier()
	port := NextNumber(8080)
	server.WithPort(port)
	server.WithNotifier(notifier)
	server.WithHostName("127.0.0.1")
	server.WithEmbeddedFileSystem(&embeddedFileSystem)
	server.WithPageController(CsrController{})
	go server.Start()
	defer server.Stop()
	time.Sleep(1 * time.Second)

	expected := "<script type=\"application/javascript\">function target(){return document.getElementById("
	actual, getError := HttpGet(fmt.Sprintf("http://127.0.0.1:%d/", port), nil)
	if getError != nil {
		test.Fatal(getError)
	}

	ok := strings.Contains(actual, expected)

	if !ok {
		test.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}
