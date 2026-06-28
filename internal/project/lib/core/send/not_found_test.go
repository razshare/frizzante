package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestNotFound(t *testing.T) {
	http := mocks.NewScope()
	NotFound(http, "not found")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 404 {
		t.Fatal("status should be 404")
	}
	if string(writer.MockBytes) != "not found" {
		t.Fatal("content should be not found")
	}
}
