package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestConnectionSendStatus(test *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	expected := 201
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestConnectionSendStatus", port))
	if getError != nil {
		test.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	actual := response.StatusCode

	if actual != expected {
		test.Fatalf("server was expected to respond with status code '%d', received '%d' intead", expected, actual)
	}
}

func TestConnectionSendHeader(test *testing.T) {
	lock := <-server
	defer func() { server <- lock }()

	expected := "application/json"
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestConnectionSendHeader", port))
	if getError != nil {
		test.Fatal(getError)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			log.Fatal(closeError)
		}
	}(response.Body)

	actual := response.Header.Get("Content-Type")

	if actual != expected {
		test.Fatalf("server was expected to respond with header content type '%s', received '%s' intead", expected, actual)
	}
}
