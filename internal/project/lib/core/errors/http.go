package errors

import "net/http"

// StatusCode returns the appropriate HTTP status code for the error category.
func StatusCode(err *Error) int {
	if err == nil {
		return http.StatusInternalServerError
	}
	switch err.Category {
	case CategoryValidation, CategoryInput:
		return http.StatusBadRequest
	case CategoryNotFound:
		return http.StatusNotFound
	case CategoryPermission:
		return http.StatusForbidden
	case CategoryConflict:
		return http.StatusConflict
	case CategoryRateLimit:
		return http.StatusTooManyRequests
	case CategoryTimeout:
		return http.StatusGatewayTimeout
	case CategoryConfiguration, CategoryDependency, CategoryInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}