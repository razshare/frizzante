package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestHeader(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("X-Header", "value")
	if Header(http, "X-Header") != "value" {
		t.Fatal("header should be value")
	}
}
