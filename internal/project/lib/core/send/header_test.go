package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestHeader(t *testing.T) {
	client := mocks.NewClient()
	Header(client, "key", "value")
	writer := client.Writer.(*mocks.ResponseWriter)
	if writer.MockHeader.Get("key") != "value" {
		t.Fatal("key should be value")
	}
}
