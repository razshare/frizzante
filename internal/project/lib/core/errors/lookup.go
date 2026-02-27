package errors

import "fmt"

func Missing(resource string, identifier string) *Error {
	return WithDetail(New(
		"not_found",
		CategoryNotFound,
		fmt.Sprintf("%s not found", resource),
	), fmt.Sprintf("looked for %s with identifier %q", resource, identifier))
}

func AlreadyExists(resource string, identifier string) *Error {
	return WithDetail(New(
		"already_exists",
		CategoryConflict,
		fmt.Sprintf("%s already exists", resource),
	), fmt.Sprintf("%s %q already exists", resource, identifier))
}
