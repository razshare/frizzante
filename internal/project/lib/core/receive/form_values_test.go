package receive

import (
	"bytes"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestFormValues(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="[]items"`),
			[]byte(``),
			[]byte(`apple`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="[]items"`),
			[]byte(``),
			[]byte(`banana`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="[]items"`),
			[]byte(``),
			[]byte(`cherry`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	values := FormValues(client, "[]items")
	if len(values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(values))
	}
	if values[0] != "apple" {
		t.Fatalf("expected first value to be 'apple', got '%s'", values[0])
	}
	if values[1] != "banana" {
		t.Fatalf("expected second value to be 'banana', got '%s'", values[1])
	}
	if values[2] != "cherry" {
		t.Fatalf("expected third value to be 'cherry', got '%s'", values[2])
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	values = FormValues(client, "[]items")
	if len(values) != 0 {
		t.Fatalf("expected 0 values, got %d", len(values))
	}
}
