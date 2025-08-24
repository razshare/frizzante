package receive

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestCookie(t *testing.T) {
	c := mock.NewClient()
	c.Request.Header.Set("Cookie", "cookie=monster;")
	ck := Cookie(c, "cookie")
	if ck != "monster" {
		t.Fatal("cookie should be monster")
	}
}
