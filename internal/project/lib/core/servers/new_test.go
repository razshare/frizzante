package servers

import (
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	server := New(
		log.New(os.Stderr, "[error]: ", log.Ldate|log.Ltime),
		log.New(os.Stdout, "[info]: ", log.Ldate|log.Ltime),
	)
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
