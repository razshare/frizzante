package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestContentType(t *testing.T) {
	client := mocks.NewClient()
	ContentType(client, "text/html")
	writer := client.Writer.(*mocks.ResponseWriter)
	if writer.MockHeader.Get("Content-Type") != "text/html" {
		t.Fatal("content type should be text/html")
	}
}
