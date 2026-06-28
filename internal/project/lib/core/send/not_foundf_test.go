package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestNotFoundf(t *testing.T) {
	http := mocks.NewScope()
	NotFoundf(http, "not %s", "found")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 404 {
		t.Fatal("status should be 404")
	}
	if string(writer.MockBytes) != "not found" {
		t.Fatal("content should be not found")
	}
}
