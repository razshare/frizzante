package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestNavigatef(t *testing.T) {
	client := mocks.NewClient()
	Navigatef(client, "/%s", "about")
	writer := client.Writer.(*mocks.ResponseWriter)
	if client.Status != 302 {
		t.Fatal("status should be 302")
	}
	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}
