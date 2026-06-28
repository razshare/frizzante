package receive

import (
	"bytes"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestFormParsedValueString(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key string `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != "value" {
		t.Fatal("key should be value")
	}
}

func TestFormParsedValueBool(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key bool `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || !form.Key {
		t.Fatal("key should be true")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || !form.Key {
		t.Fatal("key should be true")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || !form.Key {
		t.Fatal("key should be true")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || !form.Key {
		t.Fatal("key should be true")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || !form.Key {
		t.Fatal("key should be true")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || !form.Key {
		t.Fatal("key should be true")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = false
	if ok := Form(http, &form); ok || form.Key {
		t.Fatal("key should be false")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key {
		t.Fatal("key should be false")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key {
		t.Fatal("key should be false")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
	boundary.MockBuffer = bytes.Join(
		[][]byte{
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW`),
			[]byte(`Content-Disposition: form-data; name="key"`),
			[]byte(``),
			[]byte(`Payload`),
			[]byte(`------WebKitFormBoundary7MA4YWxkTrZu0gW--`),
		},
		[]byte("\n"),
	)
	if ok := Form(http, &form); ok || form.Key {
		t.Fatal("key should be false")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key {
		t.Fatal("key should be false")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key {
		t.Fatal("key should be false")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key {
		t.Fatal("key should be false")
	}
}

func TestFormParsedValueBools(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key []bool `form:"[]flags"`
	}
	var form FormPayload
	if !Form(http, &form) {
		t.Fatal("FormBoolSlice should return true")
	}
	if len(form.Key) != 3 {
		t.Fatalf("expected 3 values, got %d", len(form.Key))
	}
	if form.Key[0] != true {
		t.Fatal("expected first value to be true")
	}
	if form.Key[1] != false {
		t.Fatal("expected second value to be false")
	}
	if form.Key[2] != true {
		t.Fatal("expected third value to be true")
	}
}

func TestFormParsedValueFloat32(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key float32 `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 3.14 {
		t.Fatal("key should be a valid float32 with value 3.14")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key != -1.1 {
		t.Fatal("key should be a valid float32 with value -1.1")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid float32")
	}
}

func TestFormParsedValueFloat64(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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

	type FormPayload struct {
		Key float64 `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 3.14 {
		t.Fatal("key should be a valid float64 with value 3.14")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key != -1.1 {
		t.Fatal("key should be a valid float64 with value -1.1")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid float64")
	}
}

func TestFormParsedValueInt(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key int `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 5 {
		t.Fatal("key should be a valid int with value 5")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key != -11 {
		t.Fatal("key should be a valid int with value -11")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid int")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid int")
	}
}

func TestFormParsedValueInt32(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key int32 `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 5 {
		t.Fatal("key should be a valid int32 with value 5")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key != -11 {
		t.Fatal("key should be a valid int32 with value -11")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid int32")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid int32")
	}
}

func TestFormParsedValueInt64(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key int64 `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 5 {
		t.Fatal("key should be a valid int64 with value 5")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); !ok || form.Key != -11 {
		t.Fatal("key should be a valid int64 with value -11")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid int64")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid int64")
	}
}

func TestFormParsedValueUint(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key uint `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 5 {
		t.Fatal("key should be a valid int with value 5")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint")
	}
}

func TestFormParsedValueUint32(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key uint32 `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 5 {
		t.Fatal("key should be a valid int32 with value 5")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint32")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint32")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint32")
	}
}

func TestFormParsedValueUint64(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary := http.Request.Body.(*mocks.RequestBody)
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
	type FormPayload struct {
		Key uint64 `form:"key"`
	}
	var form FormPayload
	if ok := Form(http, &form); !ok || form.Key != 5 {
		t.Fatal("key should be a valid int64 with value 5")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint64")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint64")
	}
	http = mocks.NewScope()
	http.Request.Header.Set("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	boundary = http.Request.Body.(*mocks.RequestBody)
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
	form.Key = 0
	if ok := Form(http, &form); ok || form.Key != 0 {
		t.Fatal("key should not be a valid uint64")
	}
}
