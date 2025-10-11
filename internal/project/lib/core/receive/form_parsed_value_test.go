package receive

import (
	"bytes"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestFormParsedValueString(t *testing.T) {
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

	var value string
	if ok := FormParsedValue(client, "key", &value); !ok || value != "value" {
		t.Fatal("key should be value")
	}
}

func TestFormParsedValueBool(t *testing.T) {
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
	if ok := FormParsedValue(client, "key", &value); !ok || !value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || !value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || !value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || !value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || !value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || !value {
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

	if ok := FormParsedValue(client, "key", &value); ok || value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value {
		t.Fatal("key should be false")
	}
}

func TestFormParsedValueBools(t *testing.T) {
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
	if !FormParsedValue(client, "[]flags", &values) {
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

func TestFormParsedValueFloat32(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`3.14`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	var value float32

	if ok := FormParsedValue(client, "key", &value); !ok || value != 3.14 {
		t.Fatal("key should be a valid float32 with value 3.14")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`-1.1`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormParsedValue(client, "key", &value); !ok || value != -1.1 {
		t.Fatal("key should be a valid float32 with value -1.1")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`1qwerty`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid float32")
	}
}

func TestFormParsedValueFloat64(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`3.14`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	var value float64

	if ok := FormParsedValue(client, "key", &value); !ok || value != 3.14 {
		t.Fatal("key should be a valid float64 with value 3.14")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`-1.1`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormParsedValue(client, "key", &value); !ok || value != -1.1 {
		t.Fatal("key should be a valid float64 with value -1.1")
	}

	client = mock.NewClient()
	client.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = client.Request.Body.(*mock.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`1qwerty`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid float64")
	}
}

func TestFormParsedValueInt(t *testing.T) {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value != 5 {
		t.Fatal("key should be a valid int with value 5")
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

	if ok := FormParsedValue(client, "key", &value); !ok || value != -11 {
		t.Fatal("key should be a valid int with value -11")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid int")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid int")
	}
}

func TestFormParsedValueInt32(t *testing.T) {
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

	var value int32

	if ok := FormParsedValue(client, "key", &value); !ok || value != 5 {
		t.Fatal("key should be a valid int32 with value 5")
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

	if ok := FormParsedValue(client, "key", &value); !ok || value != -11 {
		t.Fatal("key should be a valid int32 with value -11")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid int32")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid int32")
	}
}

func TestFormParsedValueInt64(t *testing.T) {
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

	var value int64

	if ok := FormParsedValue(client, "key", &value); !ok || value != 5 {
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

	if ok := FormParsedValue(client, "key", &value); !ok || value != -11 {
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid int64")
	}
}

func TestFormParsedValueUint(t *testing.T) {
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

	var value uint

	if ok := FormParsedValue(client, "key", &value); !ok || value != 5 {
		t.Fatal("key should be a valid int with value 5")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint")
	}
}

func TestFormParsedValueUint32(t *testing.T) {
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

	var value uint32

	if ok := FormParsedValue(client, "key", &value); !ok || value != 5 {
		t.Fatal("key should be a valid int32 with value 5")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint32")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint32")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint32")
	}
}

func TestFormParsedValueUint64(t *testing.T) {
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

	var value uint64

	if ok := FormParsedValue(client, "key", &value); !ok || value != 5 {
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint64")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint64")
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

	if ok := FormParsedValue(client, "key", &value); ok || value != 0 {
		t.Fatal("key should not be a valid uint64")
	}
}
