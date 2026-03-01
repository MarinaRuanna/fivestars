package customerror

const UnauthorizedErrorType ErrorType = "unauthorized"

func NewUnauthorizedError(message string) error {
	return &CustomError{
		messagePrefix: "Unauthorized",
		message:       message,
		statusCode:    401,
		errorType:     UnauthorizedErrorType,
	}
}

func (e CustomError) IsUnauthorizedError() bool {
	return e.ErrorType() == UnauthorizedErrorType
}
