package errors

import "errors"

// As extracts a structured *Error from err if one exists in the chain, else nil.
func As(err error) *Error {
	var target *Error
	if errors.As(err, &target) {
		return target
	}
	return nil
}

// Is reports whether err matches target using the standard errors.Is.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// Join combines multiple errors into a single error using the standard errors.Join, included to not lose stdlib error construction.
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Plain creates a plain unstructured error using the standard errors.New, included to not lose stdlib error construction.
func Plain(text string) error {
	return errors.New(text)
}