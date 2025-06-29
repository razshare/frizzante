package main

import (
	"bytes"
	"fmt"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/nums"
	"github.com/razshare/frizzante/sessions"
	"github.com/razshare/frizzante/web"
	"io"
	"net/http"
	"testing"
	"time"
)

type State struct {
	Name string
}

func TestSession(t *testing.T) {
	port := nums.NextNumber(8080)
	server := web.NewServer()
	server.ViewRoot = "templates/project"
	server.PublicRoot = "templates/project/dist/client"
	server.ViewServer = "templates/project/dist/server.js"
	server.ViewIndex = "templates/project/dist/client/index.html"
	server.Efs = tefs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(web.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		state, _ := sessions.Start(con, State{Name: "test"})
		con.SendMessage(fmt.Sprintf("hello %s", state.Name))
	}})
	server.AddRoute(web.Route{Pattern: "POST /", Handler: func(con *connections.Connection) {
		state, operator := sessions.Start(con, State{})
		defer operator.Save(state)
		state.Name = con.ReceiveMessage()
	}})

	go server.Start()
	defer func() { server.Stop() }()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer response.Body.Close()

	readBytes, _ := io.ReadAll(response.Body)
	readString := string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello test`")
	}

	response, getError = http.Post(fmt.Sprintf("http://127.0.0.1:%d/", port), "text/plain", bytes.NewBufferString("world"))
	if getError != nil {
		t.Fatal(getError)
	}
	defer response.Body.Close()

	response, getError = http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer response.Body.Close()

	readBytes, _ = io.ReadAll(response.Body)
	readString = string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello world`")
	}
}

func TestSessionExpectFail(t *testing.T) {
	port := nums.NextNumber(8080)
	server := web.NewServer()
	server.ViewRoot = "templates/project"
	server.PublicRoot = "templates/project/dist/client"
	server.ViewServer = "templates/project/dist/server.js"
	server.ViewIndex = "templates/project/dist/client/index.html"
	server.Efs = tefs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(web.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		state, _ := sessions.Start(con, State{Name: "test"})
		con.SendMessage(fmt.Sprintf("hello %s", state.Name))
	}})
	server.AddRoute(web.Route{Pattern: "POST /", Handler: func(con *connections.Connection) {
		state, _ := sessions.Start(con, State{})
		// Without this, session state should not be updated.
		//defer operator.Save(state)
		state.Name = con.ReceiveMessage()
	}})

	go server.Start()
	defer func() { server.Stop() }()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer response.Body.Close()

	readBytes, _ := io.ReadAll(response.Body)
	readString := string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello test`")
	}

	response, getError = http.Post(fmt.Sprintf("http://127.0.0.1:%d/", port), "text/plain", bytes.NewBufferString("world"))
	if getError != nil {
		t.Fatal(getError)
	}
	defer response.Body.Close()

	response, getError = http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer response.Body.Close()

	readBytes, _ = io.ReadAll(response.Body)
	readString = string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello test`")
	}
}
