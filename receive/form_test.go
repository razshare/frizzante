package receive

import (
	"bytes"
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestForm(t *testing.T) {
	c := mock.NewClient()
	c.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	b := c.Request.Body.(*mock.RequestBody)
	b.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`value`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if Form(c).Get("key") != "value" {
		t.Fatal("key should be value")
	}
}
