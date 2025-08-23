package send

import "testing"

func TestCookie(t *testing.T) {
	client := MockClient()
	Cookie(client, "cookie", "monster")
	writer := client.Writer.(*MockWriter)
	if writer.MockHeader.Get("Set-Cookie") != "cookie=monster; Path=/; HttpOnly" {
		t.Fatal("cookie should be monster")
	}
}
