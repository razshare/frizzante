package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestPath(t *testing.T) {
	http := mocks.NewScope()
	http.Request.SetPathValue("key", "value")
	if Path(http, "key") != "value" {
		t.Fatal("key should be value")
	}
}
