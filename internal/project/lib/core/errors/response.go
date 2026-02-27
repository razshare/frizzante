package errors

type Response struct {
	Code       string   `json:"code"`
	Category   Category `json:"category"`
	Message    string   `json:"message"`
	Field      string   `json:"field,omitempty"`
	Suggestion string   `json:"suggestion,omitempty"`
	Status     int      `json:"status"`
}

func ToResponse(err *Error) Response {
	if err == nil {
		return Response{
			Code:     "internal",
			Category: CategoryInternal,
			Message:  "an unknown error occurred",
			Status:   StatusCode(nil),
		}
	}
	return Response{
		Code:       err.Code,
		Category:   err.Category,
		Message:    err.Message,
		Field:      err.Field,
		Suggestion: err.Suggestion,
		Status:     StatusCode(err),
	}
}