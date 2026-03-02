package customerror

const MethodNotAllowedErrorType ErrorType = "method_not_allowed"

func NewMethodNotAllowedError(message string) error {
	return &CustomError{
		messagePrefix: "Method not allowed",
		message:       message,
		statusCode:    405,
		errorType:     MethodNotAllowedErrorType,
	}
}

func (e CustomError) IsMethodNotAllowedError() bool {
	return e.ErrorType() == MethodNotAllowedErrorType
}
