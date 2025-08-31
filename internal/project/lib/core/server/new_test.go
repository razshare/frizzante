package server

import "testing"

func TestNew(t *testing.T) {
	s := New()

	if s.InfoLog == nil {
		t.Fatal("server should have an info log")
	}

	if s.ErrorLog == nil {
		t.Fatal("server should have an error log")
	}

	if s.PublicRoot == "" {
		t.Fatal("server should have a public root")
	}

	if s.Channels.Stop == nil {
		t.Fatal("server should have a stop channel")
	}

	if s.Addr == "" {
		t.Fatal("server should have an address")
	}

	if s.Handler == nil {
		t.Fatal("server should have a mux")
	}
}
