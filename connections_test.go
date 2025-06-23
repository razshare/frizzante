package frizzante

import (
	"fmt"
	"github.com/razshare/frizzante/connections"
	"github.com/razshare/frizzante/nums"
	"github.com/razshare/frizzante/routes"
	"github.com/razshare/frizzante/server"
	ghttp "net/http"
	"testing"
	"time"
)

func TestConnection_SendStatus(test *testing.T) {
	expected := 201
	port := nums.NextNumber(8080)
	server.WithEfs(emb)
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.AddRoute(routes.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		con.SendStatus(expected)
		con.SendMessage("ok")
	}})

	go server.Start()
	defer func() {
		server.Stop()
		server.Reset()
	}()

	time.Sleep(1 * time.Second)

	response, getError := ghttp.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
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
	port := nums.NextNumber(8080)
	server.WithEfs(emb)
	server.WithAddress(fmt.Sprintf("127.0.0.1:%d", port))
	server.AddRoute(routes.Route{Pattern: "GET /", Handler: func(con *connections.Connection) {
		con.SendHeader("Content-Type", expected)
		con.SendMessage("{}")
	}})

	go server.Start()
	defer func() {
		server.Stop()
		server.Reset()
	}()

	time.Sleep(1 * time.Second)

	response, getError := ghttp.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if getError != nil {
		test.Fatal(getError)
	}
	defer response.Body.Close()

	actual := response.Header.Get("Content-Type")

	if actual != expected {
		test.Fatalf("server was expected to respond with header content type '%s', received '%s' intead", expected, actual)
	}
}
