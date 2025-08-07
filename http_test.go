package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestServerAddRoute(test *testing.T) {
	lock := <-Server
	defer func() { Server <- lock }()

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestServerAddRoute", Port))
	if getError != nil {
		test.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		test.Fatal(readAllError)
	}

	expected := "hello"
	actual := string(readAllBytes)

	if actual != expected {
		test.Fatalf("server was expected to respond with '%s', received '%s' instead", expected, actual)
	}
}
