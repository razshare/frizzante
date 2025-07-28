package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

type State struct {
	Name string
}

func TestSession(t *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSession", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	readBytes, _ := io.ReadAll(response.Body)
	readString := string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello test`")
	}

	response, getError = http.Post(fmt.Sprintf("http://127.0.0.1:%d/TestSession", port), "text/plain", bytes.NewBufferString("world"))
	if getError != nil {
		t.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	response, getError = http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSession", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	readBytes, _ = io.ReadAll(response.Body)
	readString = string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello world`")
	}
}

func TestSessionExpectFail(t *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSessionExpectFail", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	readBytes, _ := io.ReadAll(response.Body)
	readString := string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello test`")
	}

	response, getError = http.Post(fmt.Sprintf("http://127.0.0.1:%d/TestSessionExpectFail", port), "text/plain", bytes.NewBufferString("world"))
	if getError != nil {
		t.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	response, getError = http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSessionExpectFail", port))
	if getError != nil {
		t.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	readBytes, _ = io.ReadAll(response.Body)
	readString = string(readBytes)

	if "hello test" != readString {
		t.Fatal("response should've been `hello test`")
	}
}
