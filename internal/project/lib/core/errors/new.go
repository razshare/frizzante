package errors

func New(code string, category Category, message string) *Error {
	return &Error{
		Code:     code,
		Category: category,
		Message:  message,
	}
}

func Wrap(code string, category Category, message string, cause error) *Error {
	return &Error{
		Code:     code,
		Category: category,
		Message:  message,
		Cause:    cause,
	}
}

func WithField(err *Error, field string) *Error {
	out := *err
	out.Field = field
	return &out
}

func WithDetail(err *Error, detail string) *Error {
	out := *err
	out.Detail = detail
	return &out
}

func WithSuggestion(err *Error, suggestion string) *Error {
	out := *err
	out.Suggestion = suggestion
	return &out
}