package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestRedirectf(t *testing.T) {
	client := mocks.NewClient()
	Redirectf(client, 303, "/%s", "about")
	writer := client.Writer.(*mocks.ResponseWriter)
	if client.Status != 303 {
		t.Fatal("status should be 303")
	}
	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}
