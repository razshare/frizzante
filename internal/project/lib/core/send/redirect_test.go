package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestRedirect(t *testing.T) {
	client := mocks.NewClient()
	Redirect(client, "/about", 303)
	writer := client.Writer.(*mocks.ResponseWriter)
	if client.Status != 303 {
		t.Fatal("status should be 303")
	}
	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}
