package customerror

const ForbiddenErrorType ErrorType = "forbidden"

func NewForbiddenError(message string) error {
	return &CustomError{
		messagePrefix: "Forbidden",
		message:       message,
		statusCode:    403,
		errorType:     ForbiddenErrorType,
	}
}

func (e CustomError) IsForbiddenError() bool {
	return e.ErrorType() == ForbiddenErrorType
}
