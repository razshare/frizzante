package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestForbiddenf(t *testing.T) {
	http := mocks.NewScope()
	Forbiddenf(http, "%s", "forbidden")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 403 {
		t.Fatal("status should be 403")
	}
	if string(writer.MockBytes) != "forbidden" {
		t.Fatal("content should be forbidden")
	}
}
