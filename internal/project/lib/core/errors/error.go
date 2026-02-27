package errors

import "fmt"

type Error struct {
	// Code is a machine-readable identifier for the error (e.g. "missing_field").
	Code string `json:"code"`
	// Category classifies the error by domain.
	Category Category `json:"category"`
	// Message is a user-safe, human-readable description of what went wrong.
	Message string `json:"message"`
	// Detail provides additional diagnostic context not shown to end users.
	Detail string `json:"detail,omitempty"`
	// Field identifies the specific input field or parameter that caused the error.
	Field string `json:"field,omitempty"`
	// Suggestion provides an actionable hint for resolving the error.
	Suggestion string `json:"suggestion,omitempty"`
	// Cause is the underlying error, if any.
	Cause error `json:"-"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func UserMessage(err *Error) string {
	if err == nil {
		return ""
	}
	if err.Suggestion != "" {
		return fmt.Sprintf("%s — %s", err.Message, err.Suggestion)
	}
	return err.Message
}