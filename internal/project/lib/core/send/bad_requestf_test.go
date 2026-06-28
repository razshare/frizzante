package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestBadRequestf(t *testing.T) {
	http := mocks.NewScope()
	BadRequestf(http, "bad %s", "request")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 400 {
		t.Fatal("status should be 400")
	}
	if string(writer.MockBytes) != "bad request" {
		t.Fatal("content should be bad request")
	}
}
