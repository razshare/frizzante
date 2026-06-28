package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestCookie(t *testing.T) {
	http := mocks.NewScope()
	Cookie(http, "cookie", "monster")
	writer := http.Writer.(*mocks.ResponseWriter)
	if writer.MockHeader.Get("Set-Cookie") != "cookie=monster; Path=/; HttpOnly" {
		t.Fatal("cookie should be monster")
	}
}
