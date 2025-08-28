package send

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestCookie(t *testing.T) {
	client := mock.NewClient()
	Cookie(client, "cookie", "monster")
	writer := client.Writer.(*mock.ResponseWriter)
	if writer.MockHeader.Get("Set-Cookie") != "cookie=monster; Path=/; HttpOnly" {
		t.Fatal("cookie should be monster")
	}
}
