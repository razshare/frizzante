package errors

import "fmt"

func Required(field string, value string) *Error {
	if value == "" {
		return WithField(New(
			"required_field",
			CategoryValidation,
			fmt.Sprintf("%s is required", field),
		), field)
	}
	return nil
}

func NotEmpty(field string, length int) *Error {
	if length == 0 {
		return WithField(New(
			"empty_collection",
			CategoryValidation,
			fmt.Sprintf("%s must not be empty", field),
		), field)
	}
	return nil
}

func InBounds(field string, index int, length int) *Error {
	if length < 0 {
		return WithField(WithDetail(New(
			"invalid_length",
			CategoryValidation,
			fmt.Sprintf("%s has an invalid length", field),
		), fmt.Sprintf("length %d is negative", length)), field)
	}
	if index < 0 || index >= length {
		return WithField(WithDetail(New(
			"index_out_of_bounds",
			CategoryValidation,
			fmt.Sprintf("%s is out of bounds", field),
		), fmt.Sprintf("index %d is outside range [0, %d)", index, length)), field)
	}
	return nil
}

