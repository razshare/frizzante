package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestCookie(t *testing.T) {
	_, writer := mocks.NewExchange()
	Cookie(writer, "cookie", "monster")
	if writer.MockHeader.Get("Set-Cookie") != "cookie=monster; Path=/; HttpOnly" {
		t.Fatal("cookie should be monster")
	}
}
