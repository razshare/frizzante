package ssr

import "testing"

func TestNewServer(t *testing.T) {
	server := NewServer()

	if server.InfoLog == nil {
		t.Fatal("server should have an info log")
	}

	if server.ErrorLog == nil {
		t.Fatal("server should have an error log")
	}

	if server.Addr == "" {
		t.Fatal("server should have an address")
	}

	if server.Handler == nil {
		t.Fatal("server should have a mux")
	}
}
