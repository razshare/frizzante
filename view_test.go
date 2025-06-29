package main

import (
	"fmt"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/nums"
	"github.com/razshare/frizzante/views"
	"github.com/razshare/frizzante/web"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRenderServer(test *testing.T) {
	port := nums.NextNumber(8080)
	server := web.NewServer()
	server.Efs = testEfs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(web.Route{Pattern: "GET /welcome", Handler: func(con *connections.Connection) {
		con.SendView(views.View{
			Name:       "Welcome",
			RenderMode: views.RenderModeServer,
			Data:       map[string]any{"name": "world"},
		})
	}})

	go server.Start()
	defer func() { server.Stop() }()
	time.Sleep(1 * time.Second)

	expected := "<h1>Welcome to Frizzante.</h1>"
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/welcome", port))
	if getError != nil {
		test.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	ok := strings.Contains(actual, expected)

	if !ok {
		test.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}

func TestRenderClient(test *testing.T) {
	port := nums.NextNumber(8080)
	server := web.NewServer()
	server.Efs = testEfs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(web.Route{Pattern: "GET /welcome", Handler: func(con *connections.Connection) {
		con.SendView(views.View{
			Name:       "Welcome",
			RenderMode: views.RenderModeClient,
			Data:       map[string]any{"name": "world"},
		})
	}})
	go server.Start()
	defer func() { server.Stop() }()

	time.Sleep(1 * time.Second)

	expected := "<script type=\"application/javascript\">function target(){return document.getElementById("
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/welcome", port))
	if getError != nil {
		test.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	ok := strings.Contains(actual, expected)

	if !ok {
		test.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}
