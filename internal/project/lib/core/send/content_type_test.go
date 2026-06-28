package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestContentType(t *testing.T) {
	http := mocks.NewScope()
	ContentType(http, "text/html")
	writer := http.Writer.(*mocks.ResponseWriter)
	if writer.MockHeader.Get("Content-Type") != "text/html" {
		t.Fatal("content type should be text/html")
	}
}
