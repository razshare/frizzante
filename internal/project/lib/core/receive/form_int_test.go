package receive

import (
	"bytes"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestFormInt(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`5`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	var value int
	var value32 int32
	var value64 int64

	if ok := FormInt(client, "key", &value); !ok || value != 5 {
		t.Fatal("key should be a valid int with value 5")
	}

	if ok := FormInt32(client, "key", &value32); !ok || value32 != 5 {
		t.Fatal("key should be a valid int32 with value 5")
	}

	if ok := FormInt64(client, "key", &value64); !ok || value64 != 5 {
		t.Fatal("key should be a valid int64 with value 5")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`-11`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormInt(client, "key", &value); !ok || value != -11 {
		t.Fatal("key should be a valid int with value -11")
	}

	if ok := FormInt32(client, "key", &value32); !ok || value32 != -11 {
		t.Fatal("key should be a valid int32 with value -11")
	}

	if ok := FormInt64(client, "key", &value64); !ok || value64 != -11 {
		t.Fatal("key should be a valid int64 with value -11")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`grweqrqa`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormInt(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid int")
	}

	if ok := FormInt32(client, "key", &value32); ok || value32 != 0 {
		t.Fatal("key should not be a valid int32")
	}

	if ok := FormInt64(client, "key", &value64); ok || value64 != 0 {
		t.Fatal("key should not be a valid int64")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`7.7`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormInt(client, "key", &value); ok || value32 != 0 {
		t.Fatal("key should not be a valid int")
	}

	if ok := FormInt32(client, "key", &value32); ok || value32 != 0 {
		t.Fatal("key should not be a valid int32")
	}

	if ok := FormInt64(client, "key", &value64); ok || value64 != 0 {
		t.Fatal("key should not be a valid int64")
	}
}
