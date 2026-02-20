package customerror

import "net/http"

const ConflictErrorType ErrorType = "conflict"

func NewConflictError(message string) error {
	return &CustomError{
		messagePrefix: "Conflict",
		message:       message,
		statusCode:    http.StatusConflict,
		errorType:     ConflictErrorType,
	}
}
