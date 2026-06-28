package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestAccept(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Accept", "text/html")
	if Accept(http) != "text/html" {
		t.Fatal("accept should be text/html")
	}
}
