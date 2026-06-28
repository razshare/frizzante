package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestUnauthorized(t *testing.T) {
	http := mocks.NewScope()
	Unauthorized(http, "unauthorized")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 401 {
		t.Fatal("status should be 401")
	}
	if string(writer.MockBytes) != "unauthorized" {
		t.Fatal("content should be unauthorized")
	}
}
