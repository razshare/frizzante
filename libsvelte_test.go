package frizzante

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRenderServer(test *testing.T) {
	server := ServerCreate()
	notifier := NotifierCreate()
	port := NextNumber(8080)
	ServerWithPort(server, port)
	ServerWithHostName(server, "127.0.0.1")
	ServerWithNotifier(server, notifier)
	ServerWithEmbeddedFileSystem(server, embeddedFileSystem)
	ServerWithPage(server, func(
		withPath func(path string),
		withDocument func(document *Document),
		show func(showFunction func(request *Request, response *Response, document *Document)),
		action func(actionFunction func(request *Request, response *Response, document *Document)),
	) {
		withPath("/")
		withDocument(DocumentCreate("Welcome"))
		show(func(request *Request, response *Response, doc *Document) {
			doc.Render = RenderServer
			doc.Data["name"] = "world"
		})
	})
	go ServerStart(server)
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
	server := ServerCreate()
	notifier := NotifierCreate()
	port := NextNumber(8080)
	ServerWithPort(server, port)
	ServerWithNotifier(server, notifier)
	ServerWithHostName(server, "127.0.0.1")
	ServerWithEmbeddedFileSystem(server, embeddedFileSystem)
	ServerWithPage(server, func(
		withPath func(path string),
		withDocument func(document *Document),
		withBaseHandler func(baseHandler func(request *Request, response *Response, doc *Document)),
		withActionHandler func(actionFunction func(request *Request, response *Response, doc *Document)),
	) {
		withPath("/")
		withDocument(DocumentCreate("Welcome"))
		withBaseHandler(func(request *Request, response *Response, doc *Document) {
			doc.Render = RenderClient
			doc.Data["name"] = "world"
		})
	})
	go ServerStart(server)
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
