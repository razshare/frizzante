package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestNavigate(t *testing.T) {
	_, writer := mocks.NewExchange()
	if err := Navigate(writer, "/about"); err != nil {
		t.Fatal("navigating should succeed")
	}
	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}
