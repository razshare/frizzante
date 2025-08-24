package receive

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestHeader(t *testing.T) {
	c := mock.NewClient()
	c.Request.Header.Set("X-Header", "value")
	if Header(c, "X-Header") != "value" {
		t.Fatal("header should be value")
	}
}

func TestContentType(t *testing.T) {
	c := mock.NewClient()
	c.Request.Header.Set("Content-Type", "text/html")
	if ContentType(c) != "text/html" {
		t.Fatal("content type should be text/html")
	}
}

func TestAccept(t *testing.T) {
	c := mock.NewClient()
	c.Request.Header.Set("Accept", "text/html")
	if Accept(c) != "text/html" {
		t.Fatal("accept should be text/html")
	}
}
