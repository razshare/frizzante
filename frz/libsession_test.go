package frz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

type State struct {
	Name string
}

func TestSession(t *testing.T) {
	port := NextNumber(8080)
	server := NewServer().
		WithEfs(dist).
		WithAddress(fmt.Sprintf("127.0.0.1:%d", port)).
		AddRoute(Route{Pattern: "GET /", Handler: func(c *Connection) {
			state, _ := Session(c, State{Name: "test"})
			c.SendMessage(fmt.Sprintf("hello %s", state.Name))
		}}).
		AddRoute(Route{Pattern: "POST /", Handler: func(c *Connection) {
			state, operator := Session(c, State{})
			defer operator.Save(state)
			state.Name = c.ReceiveMessage()
		}})

	go server.Start()
	defer server.Stop()

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
	port := NextNumber(8080)
	server := NewServer().
		WithEfs(dist).
		WithAddress(fmt.Sprintf("127.0.0.1:%d", port)).
		AddRoute(Route{Pattern: "GET /", Handler: func(c *Connection) {
			state, _ := Session(c, State{Name: "test"})
			c.SendMessage(fmt.Sprintf("hello %s", state.Name))
		}}).
		AddRoute(Route{Pattern: "POST /", Handler: func(c *Connection) {
			state, _ := Session(c, State{})
			// Without this, session state should not be updated.
			//defer operator.Save(state)
			state.Name = c.ReceiveMessage()
		}})

	go server.Start()
	defer server.Stop()

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
