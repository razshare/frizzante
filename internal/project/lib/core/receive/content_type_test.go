package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestContentType(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "text/html")
	if ContentType(http) != "text/html" {
		t.Fatal("content type should be text/html")
	}
}
