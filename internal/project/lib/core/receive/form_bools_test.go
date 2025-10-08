package receive

import (
	"bytes"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestFormBools(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="[]flags"`),
			[]byte(``),
			[]byte(`true`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="[]flags"`),
			[]byte(``),
			[]byte(`false`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="[]flags"`),
			[]byte(``),
			[]byte(`1`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	var values []bool
	if !FormBools(client, "[]flags", &values) {
		t.Fatal("FormBoolSlice should return true")
	}
	if len(values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(values))
	}
	if values[0] != true {
		t.Fatal("expected first value to be true")
	}
	if values[1] != false {
		t.Fatal("expected second value to be false")
	}
	if values[2] != true {
		t.Fatal("expected third value to be true")
	}
}
