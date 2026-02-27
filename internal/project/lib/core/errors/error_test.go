package errors

import (
	"errors"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	err := New("test_code", CategoryValidation, "field is required")
	if err.Error() != "field is required" {
		t.Fatalf("expected 'field is required', got '%s'", err.Error())
	}
}

func TestErrorMessageWithCause(t *testing.T) {
	cause := errors.New("underlying issue")
	err := Wrap("test_code", CategoryInternal, "operation failed", cause)
	expected := "operation failed: underlying issue"
	if err.Error() != expected {
		t.Fatalf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestUnwrap(t *testing.T) {
	cause := errors.New("underlying issue")
	err := Wrap("test_code", CategoryInternal, "operation failed", cause)
	if !errors.Is(err, cause) {
		t.Fatal("errors.Is should match the underlying cause")
	}
}

func TestAsStructured(t *testing.T) {
	err := New("test_code", CategoryValidation, "bad input")
	var wrapped error = err
	result := As(wrapped)
	if result == nil {
		t.Fatal("As should extract structured error")
	}
	if result.Code != "test_code" {
		t.Fatalf("expected code 'test_code', got '%s'", result.Code)
	}
}

func TestAsPlainError(t *testing.T) {
	err := errors.New("plain error")
	result := As(err)
	if result != nil {
		t.Fatal("As should return nil for plain errors")
	}
}

func TestWithField(t *testing.T) {
	err := New("test_code", CategoryValidation, "bad input")
	result := WithField(err, "email")
	if result.Field != "email" {
		t.Fatalf("expected field 'email', got '%s'", result.Field)
	}
	if err.Field != "" {
		t.Fatal("original error should not be modified")
	}
}

func TestWithDetail(t *testing.T) {
	err := New("test_code", CategoryValidation, "bad input")
	result := WithDetail(err, "expected string got int")
	if result.Detail != "expected string got int" {
		t.Fatalf("expected detail 'expected string got int', got '%s'", result.Detail)
	}
	if err.Detail != "" {
		t.Fatal("original error should not be modified")
	}
}

func TestWithSuggestion(t *testing.T) {
	err := New("test_code", CategoryValidation, "bad input")
	result := WithSuggestion(err, "provide a valid email address")
	if result.Suggestion != "provide a valid email address" {
		t.Fatalf("expected suggestion, got '%s'", result.Suggestion)
	}
	if err.Suggestion != "" {
		t.Fatal("original error should not be modified")
	}
}

func TestUserMessage(t *testing.T) {
	err := New("test_code", CategoryValidation, "field is required")
	if UserMessage(err) != "field is required" {
		t.Fatalf("expected 'field is required', got '%s'", UserMessage(err))
	}
}

func TestUserMessageWithSuggestion(t *testing.T) {
	err := WithSuggestion(
		New("test_code", CategoryValidation, "field is required"),
		"provide a value",
	)
	expected := "field is required — provide a value"
	if UserMessage(err) != expected {
		t.Fatalf("expected '%s', got '%s'", expected, UserMessage(err))
	}
}

func TestStatusCode(t *testing.T) {
	tests := []struct {
		Category Category
		Expected int
	}{
		{CategoryValidation, 400},
		{CategoryInput, 400},
		{CategoryNotFound, 404},
		{CategoryPermission, 403},
		{CategoryConflict, 409},
		{CategoryRateLimit, 429},
		{CategoryTimeout, 504},
		{CategoryInternal, 500},
		{CategoryConfiguration, 500},
		{CategoryDependency, 500},
	}
	for _, test := range tests {
		err := New("test", test.Category, "test")
		if code := StatusCode(err); code != test.Expected {
			t.Fatalf("category %s: expected status %d, got %d", test.Category, test.Expected, code)
		}
	}
}

func TestToResponse(t *testing.T) {
	err := WithSuggestion(
		WithField(
			New("missing_field", CategoryValidation, "name is required"),
			"name",
		),
		"provide a name",
	)
	response := ToResponse(err)
	if response.Code != "missing_field" {
		t.Fatalf("expected code 'missing_field', got '%s'", response.Code)
	}
	if response.Category != CategoryValidation {
		t.Fatalf("expected category 'validation', got '%s'", response.Category)
	}
	if response.Message != "name is required" {
		t.Fatalf("expected message 'name is required', got '%s'", response.Message)
	}
	if response.Field != "name" {
		t.Fatalf("expected field 'name', got '%s'", response.Field)
	}
	if response.Suggestion != "provide a name" {
		t.Fatalf("expected suggestion 'provide a name', got '%s'", response.Suggestion)
	}
	if response.Status != 400 {
		t.Fatalf("expected status 400, got %d", response.Status)
	}
}

func TestValidateRequired(t *testing.T) {
	if err := Required("name", ""); err == nil {
		t.Fatal("Required should return error for empty string")
	}
	if err := Required("name", "Alice"); err != nil {
		t.Fatal("Required should return nil for non-empty string")
	}
}

func TestValidateInBounds(t *testing.T) {
	if err := InBounds("index", 5, 3); err == nil {
		t.Fatal("InBounds should return error when index >= length")
	}
	if err := InBounds("index", -1, 3); err == nil {
		t.Fatal("InBounds should return error when index < 0")
	}
	if err := InBounds("index", 2, 3); err != nil {
		t.Fatal("InBounds should return nil for valid index")
	}
}

func TestValidateNotEmpty(t *testing.T) {
	if err := NotEmpty("items", 0); err == nil {
		t.Fatal("NotEmpty should return error for zero length")
	}
	if err := NotEmpty("items", 3); err != nil {
		t.Fatal("NotEmpty should return nil for non-zero length")
	}
}

func TestMissing(t *testing.T) {
	err := Missing("file", "/tmp/foo.txt")
	if err.Code != "not_found" {
		t.Fatalf("expected code 'not_found', got '%s'", err.Code)
	}
	if err.Category != CategoryNotFound {
		t.Fatalf("expected category 'not_found', got '%s'", err.Category)
	}
	if err.Message != "file not found" {
		t.Fatalf("expected message 'file not found', got '%s'", err.Message)
	}
	if err.Detail == "" {
		t.Fatal("expected detail to be set")
	}
}

func TestAlreadyExists(t *testing.T) {
	err := AlreadyExists("directory", "lib/core")
	if err.Code != "already_exists" {
		t.Fatalf("expected code 'already_exists', got '%s'", err.Code)
	}
	if err.Category != CategoryConflict {
		t.Fatalf("expected category 'conflict', got '%s'", err.Category)
	}
	if err.Message != "directory already exists" {
		t.Fatalf("expected message 'directory already exists', got '%s'", err.Message)
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("disk full")
	err := Wrap("write_failed", CategoryInternal, "could not save file", cause)
	if err.Cause != cause {
		t.Fatal("Cause should be the original error")
	}
	if !errors.Is(err, cause) {
		t.Fatal("Unwrap chain should expose the cause")
	}
}

func TestWithFieldDoesNotMutateCause(t *testing.T) {
	original := New("test", CategoryValidation, "bad")
	original.Cause = errors.New("root cause")
	copy := WithField(original, "name")
	if copy.Cause != original.Cause {
		t.Fatal("WithField should preserve the Cause")
	}
}

func TestStatusCodeUnknownCategory(t *testing.T) {
	err := New("test", Category("unknown"), "test")
	if code := StatusCode(err); code != 500 {
		t.Fatalf("expected status 500 for unknown category, got %d", code)
	}
}

func TestStatusCodeNil(t *testing.T) {
	if code := StatusCode(nil); code != 500 {
		t.Fatalf("expected status 500 for nil error, got %d", code)
	}
}

func TestUserMessageNil(t *testing.T) {
	if msg := UserMessage(nil); msg != "" {
		t.Fatalf("expected empty string for nil error, got '%s'", msg)
	}
}

func TestToResponseNil(t *testing.T) {
	response := ToResponse(nil)
	if response.Code != "internal" {
		t.Fatalf("expected code 'internal', got '%s'", response.Code)
	}
	if response.Category != CategoryInternal {
		t.Fatalf("expected category 'internal', got '%s'", response.Category)
	}
	if response.Status != 500 {
		t.Fatalf("expected status 500, got %d", response.Status)
	}
}

func TestToResponseDerivesStatusCode(t *testing.T) {
	err := New("not_found", CategoryNotFound, "resource not found")
	response := ToResponse(err)
	if response.Status != 404 {
		t.Fatalf("expected status 404, got %d", response.Status)
	}
}

func TestInBoundsZeroLength(t *testing.T) {
	if err := InBounds("index", 0, 0); err == nil {
		t.Fatal("InBounds should return error when length is 0")
	}
}

func TestInBoundsNegativeLength(t *testing.T) {
	err := InBounds("index", 0, -1)
	if err == nil {
		t.Fatal("InBounds should return error when length is negative")
	}
	if err.Code != "invalid_length" {
		t.Fatalf("expected code 'invalid_length', got '%s'", err.Code)
	}
}

func TestJoin(t *testing.T) {
	err1 := errors.New("first")
	err2 := errors.New("second")
	joined := Join(err1, err2)
	if joined == nil {
		t.Fatal("Join should return a non-nil error")
	}
}

func TestPlain(t *testing.T) {
	err := Plain("something went wrong")
	if err == nil {
		t.Fatal("Plain should return a non-nil error")
	}
	if err.Error() != "something went wrong" {
		t.Fatalf("expected 'something went wrong', got '%s'", err.Error())
	}
}