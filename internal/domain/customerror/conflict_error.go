package customerror

const ConflictErrorType ErrorType = "conflict"

func NewConflictError(message string) error {
	return &CustomError{
		messagePrefix: "Conflict",
		message:       message,
		statusCode:    409,
		errorType:     ConflictErrorType,
	}
}

func (e CustomError) IsConflictError() bool {
	return e.ErrorType() == ConflictErrorType
}
