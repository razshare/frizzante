package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestRedirect(t *testing.T) {
	_, writer := mocks.NewExchange()
	if err := Redirect(writer, "/about", 303); err != nil {
		t.Fatal("redirect should succeed")
	}
	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}
