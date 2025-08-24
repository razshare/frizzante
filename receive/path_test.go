package receive

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestPath(t *testing.T) {
	c := mock.NewClient()
	c.Request.SetPathValue("key", "value")
	if Path(c, "key") != "value" {
		t.Fatal("key should be value")
	}
}
