package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestServerAddRoute(t *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestServerAddRoute", port))
	if getError != nil {
		t.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		t.Fatal(readAllError)
	}

	expected := "hello"
	actual := string(readAllBytes)

	if actual != expected {
		t.Fatalf("server was expected to respond with '%s', received '%s' instead", expected, actual)
	}
}
