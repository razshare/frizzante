package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestNavigate(t *testing.T) {
	http := mocks.NewScope()
	Navigate(http, "/about")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 302 {
		t.Fatal("status should be 302")
	}
	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}
