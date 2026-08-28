package errors

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationError(field string, message string) ValidationError {
	return ValidationError{Field: field, Message: message}
}

func (e ValidationError) Error() string {
	return e.Message
}
