package receive

import (
	"bytes"
	"main/lib/core/mock"
	"testing"
)

func TestForm(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`value`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if Form(client).Get("key") != "value" {
		t.Fatal("key should be value")
	}
}
