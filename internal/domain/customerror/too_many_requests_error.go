package customerror

const TooManyRequestsErrorType ErrorType = "too_many_requests"

func NewTooManyRequestsError(message string) error {
	return &CustomError{
		messagePrefix: "Too many requests",
		message:       message,
		statusCode:    429,
		errorType:     TooManyRequestsErrorType,
	}
}

func (e CustomError) IsTooManyRequestsError() bool {
	return e.ErrorType() == TooManyRequestsErrorType
}
