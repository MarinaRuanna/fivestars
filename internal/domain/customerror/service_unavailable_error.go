package customerror

const ServiceUnavailableErrorType ErrorType = "service_unavailable"

func NewServiceUnavailableError(message string) error {
	return &CustomError{
		messagePrefix: "Service unavailable",
		message:       message,
		statusCode:    503,
		errorType:     ServiceUnavailableErrorType,
	}
}

func (e CustomError) IsServiceUnavailableError() bool {
	return e.ErrorType() == ServiceUnavailableErrorType
}
