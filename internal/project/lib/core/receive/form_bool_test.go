package receive

import (
	"bytes"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestFormBool(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`1`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	var value bool

	if ok := FormBool(client, "key", &value); !ok || !value {
		t.Fatal("key should be true")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`t`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || !value {
		t.Fatal("key should be true")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`T`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || !value {
		t.Fatal("key should be true")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`true`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || !value {
		t.Fatal("key should be true")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`True`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || !value {
		t.Fatal("key should be true")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`TRUE`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || !value {
		t.Fatal("key should be true")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`qwerty`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); ok || value {
		t.Fatal("key should be false")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`0`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || value {
		t.Fatal("key should be false")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`f`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || value {
		t.Fatal("key should be false")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`F`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || value {
		t.Fatal("key should be false")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`false`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || value {
		t.Fatal("key should be false")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`False`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || value {
		t.Fatal("key should be false")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`FALSE`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormBool(client, "key", &value); !ok || value {
		t.Fatal("key should be false")
	}
}
