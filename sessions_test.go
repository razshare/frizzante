package main

import (
	"bytes"
	"fmt"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/nums"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/servers"
	"github.com/razshare/frizzante/sessions"
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
	server := servers.New()
	server.Efs = testEfs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(routes.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{Name: "test"}).Start()
		con.SendMessage(fmt.Sprintf("hello %s", session.State.Name))
	}})
	server.AddRoute(routes.Route{Pattern: "POST /", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{}).Start()
		defer session.Save()
		session.State.Name = con.ReceiveMessage()
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
	server := servers.New()
	server.Efs = testEfs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(routes.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{Name: "test"}).Start()
		con.SendMessage(fmt.Sprintf("hello %s", session.State.Name))
	}})
	server.AddRoute(routes.Route{Pattern: "POST /", Handler: func(con *connections.Connection) {
		session := sessions.New(con, State{}).Start()
		// Without this, session state should not be updated.
		//defer operator.Save(state)
		session.State.Name = con.ReceiveMessage()
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
