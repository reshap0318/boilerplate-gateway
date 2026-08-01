package helpers

import "errors"

// FieldError represents an error tied to a specific field.
type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return e.Message
}

// CustomError represents a custom error with HTTP status and message.
type CustomError struct {
	Status  int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// Common application errors.
var (
	ErrNotFound          = errors.New("record not found")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidCredential = errors.New("invalid email or password")
)
