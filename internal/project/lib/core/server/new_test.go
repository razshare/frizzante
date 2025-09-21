package server

import "testing"

func TestNew(t *testing.T) {
	server := New()

	if server.InfoLog == nil {
		t.Fatal("server should have an info log")
	}

	if server.ErrorLog == nil {
		t.Fatal("server should have an error log")
	}

	if server.PublicRoot == "" {
		t.Fatal("server should have a public root")
	}

	if server.Channels.Stop == nil {
		t.Fatal("server should have a stop channel")
	}

	if server.Addr == "" {
		t.Fatal("server should have an address")
	}

	if server.Handler == nil {
		t.Fatal("server should have a mux")
	}
}
