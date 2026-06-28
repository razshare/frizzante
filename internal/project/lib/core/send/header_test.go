package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestHeader(t *testing.T) {
	http := mocks.NewScope()
	Header(http, "key", "value")
	writer := http.Writer.(*mocks.ResponseWriter)
	if writer.MockHeader.Get("key") != "value" {
		t.Fatal("key should be value")
	}
}
