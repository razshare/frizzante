package frizzante

import (
	"bytes"
	"fmt"
	"github.com/razshare/frizzante/libcon"
	"github.com/razshare/frizzante/libnum"
	"github.com/razshare/frizzante/libsession"
	"github.com/razshare/frizzante/libsrv"
	"io"
	"net/http"
	"testing"
	"time"
)

type State struct {
	Name string
}

func TestSession(t *testing.T) {
	port := libnum.NextNumber(8080)
	server := libsrv.NewServer()
	server.Efs = efs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(libsrv.Route{Pattern: "GET /", Handler: func(c *libcon.Connection) {
		state, _ := libsession.Session(c, State{Name: "test"})
		c.SendMessage(fmt.Sprintf("hello %s", state.Name))
	}}).
		AddRoute(libsrv.Route{Pattern: "POST /", Handler: func(c *libcon.Connection) {
			state, operator := libsession.Session(c, State{})
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
	port := libnum.NextNumber(8080)
	server := libsrv.NewServer()
	server.Efs = efs
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(libsrv.Route{Pattern: "GET /", Handler: func(c *libcon.Connection) {
		state, _ := libsession.Session(c, State{Name: "test"})
		c.SendMessage(fmt.Sprintf("hello %s", state.Name))
	}}).
		AddRoute(libsrv.Route{Pattern: "POST /", Handler: func(c *libcon.Connection) {
			state, _ := libsession.Session(c, State{})
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
