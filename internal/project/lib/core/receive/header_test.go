package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestHeader(t *testing.T) {
	client := mocks.NewClient()
	client.Request.Header.Set("X-Header", "value")
	if Header(client, "X-Header") != "value" {
		t.Fatal("header should be value")
	}
}

func TestContentType(t *testing.T) {
	client := mocks.NewClient()
	client.Request.Header.Set("Content-Type", "text/html")
	if ContentType(client) != "text/html" {
		t.Fatal("content type should be text/html")
	}
}

func TestAccept(t *testing.T) {
	client := mocks.NewClient()
	client.Request.Header.Set("Accept", "text/html")
	if Accept(client) != "text/html" {
		t.Fatal("accept should be text/html")
	}
}
