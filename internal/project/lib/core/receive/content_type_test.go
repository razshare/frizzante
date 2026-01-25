package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestContentType(t *testing.T) {
	client := mocks.NewClient()
	client.Request.Header.Set("Content-Type", "text/html")
	if ContentType(client) != "text/html" {
		t.Fatal("content type should be text/html")
	}
}
