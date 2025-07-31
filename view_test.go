package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRenderServer(t *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	expected := "<h1>Welcome to Frizzante.</h1>"
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRenderServer", port))
	if getError != nil {
		t.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		t.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	ok := strings.Contains(actual, expected)

	if !ok {
		t.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}

func TestRenderClient(t *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	expected := "<script type=\"application/javascript\">function target(){return document.getElementById("
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRenderClient", port))
	if getError != nil {
		t.Fatal(getError)
	}

	readAllBytes, readAllError := io.ReadAll(response.Body)
	if readAllError != nil {
		t.Fatal(readAllError)
	}

	actual := string(readAllBytes)

	ok := strings.Contains(actual, expected)

	if !ok {
		t.Fatalf("server was expected to respond with a string that contains '%s', received '%s' instead", expected, actual)
	}
}
