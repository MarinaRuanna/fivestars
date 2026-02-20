package customerror

import "net/http"

const ValidationErrorType ErrorType = "validation_error"

func NewValidationError(message string) error {
	return &CustomError{
		messagePrefix: "Validation error - Message:",
		message:       message,
		statusCode:    http.StatusBadRequest,
		errorType:     ValidationErrorType,
	}
}

func (e CustomError) IsValidationError() bool {
	return e.ErrorType() == ValidationErrorType
}
