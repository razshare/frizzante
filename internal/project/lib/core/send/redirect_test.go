package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestRedirect(t *testing.T) {
	http := mocks.NewScope()
	Redirect(http, "/about", 303)
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 303 {
		t.Fatal("status should be 303")
	}
	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}
