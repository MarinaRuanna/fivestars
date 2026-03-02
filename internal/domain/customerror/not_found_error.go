package customerror

const NotFoundErrorType ErrorType = "not_found_error"

func NewNotFoundError(message string) error {
	return &CustomError{
		messagePrefix: "Not found error - Message:",
		message:       message,
		statusCode:    404,
		errorType:     NotFoundErrorType,
	}
}

func (e CustomError) IsNotFoundError() bool {
	return e.ErrorType() == NotFoundErrorType
}
