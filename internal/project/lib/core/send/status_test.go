package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestStatus(t *testing.T) {
	http := mocks.NewScope()
	Status(http, 400)
	if http.Status != 400 {
		t.Fatal("status should be 400")
	}
}
