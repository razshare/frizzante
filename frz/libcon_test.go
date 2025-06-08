package frz

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestConnection_SendStatus(test *testing.T) {
	expected := 201
	port := NextNumber(8080)
	server := NewServer().
		WithEfs(dist).
		WithAddress(fmt.Sprintf("127.0.0.1:%d", port)).
		AddRoute(Route{Pattern: "GET /", Handler: func(c *Connection) {
			c.SendStatus(expected)
			c.SendMessage("ok")
		}})

	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		test.Fatal(getError)
	}
	defer response.Body.Close()

	actual := response.StatusCode

	if actual != expected {
		test.Fatalf("server was expected to respond with status code '%d', received '%d' intead", expected, actual)
	}
}

func TestConnection_SendHeader(test *testing.T) {
	expected := "application/json"
	port := NextNumber(8080)
	server := NewServer().
		WithEfs(dist).
		WithAddress(fmt.Sprintf("127.0.0.1:%d", port)).
		AddRoute(Route{Pattern: "GET /", Handler: func(c *Connection) {
			c.SendHeader("Content-Type", expected)
			c.SendMessage("{}")
		}})

	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		test.Fatal(getError)
	}
	defer response.Body.Close()

	actual := response.Header.Get("Content-Type")

	if actual != expected {
		test.Fatalf("server was expected to respond with header content type '%s', received '%s' intead", expected, actual)
	}
}
