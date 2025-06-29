package main

import (
	"fmt"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/nums"
	"github.com/razshare/frizzante/web"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestServer_AddRoute(test *testing.T) {
	expected := "hello"
	server := web.NewServer()
	port := nums.NextNumber(8080)
	server.Address = fmt.Sprintf("127.0.0.1:%d", port)
	server.AddRoute(web.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		con.SendMessage(expected)
	}})

	go server.Start()
	defer func() { server.Stop() }()

	time.Sleep(1 * time.Second)

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		test.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	if actual != expected {
		test.Fatalf("server was expected to respond with '%s', received '%s' instead", expected, actual)
	}
}
