package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	<-serve
	defer func() { serve <- 0 }()

	r, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestRoutes", port))
	if err != nil {
		t.Fatal(err)
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	ex := "hello"
	ac := string(d)

	if ac != ex {
		t.Fatalf("server was expected to respond with '%s', received '%s' instead", ex, ac)
	}
}

func TestSendStatus(t *testing.T) {
	<-serve
	defer func() { serve <- 0 }()

	ex := 201
	r, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSendStatus", port))
	if err != nil {
		t.Fatal(err)
	}
	defer func(b io.ReadCloser) {
		err = b.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(r.Body)

	ac := r.StatusCode

	if ac != ex {
		t.Fatalf("server was expected to respond with status code '%d', received '%d' intead", ex, ac)
	}
}

func TestSendHeader(t *testing.T) {
	<-serve
	defer func() { serve <- 0 }()

	ex := "application/json"
	r, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/TestSendHeader", port))
	if err != nil {
		t.Fatal(err)
	}
	defer func(b io.ReadCloser) {
		err = b.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(r.Body)

	ac := r.Header.Get("Content-Type")

	if ac != ex {
		t.Fatalf("server was expected to respond with header content type '%s', received '%s' intead", ex, ac)
	}
}
