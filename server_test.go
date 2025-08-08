package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestRoutes(test *testing.T) {
	<-ready
	defer func() { ready <- 0 }()

	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRoutes", port))
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

func TestSendStatus(test *testing.T) {
	<-ready
	defer func() { ready <- 0 }()

	expected := 201
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSendStatus", port))
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

func TestSendHeader(test *testing.T) {
	<-ready
	defer func() { ready <- 0 }()

	expected := "application/json"
	response, getError := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSendHeader", port))
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
