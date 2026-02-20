package customerror

import "net/http"

const UnauthorizedErrorType ErrorType = "unauthorized"

func NewUnauthorizedError(message string) error {
	return &CustomError{
		messagePrefix: "Unauthorized",
		message:       message,
		statusCode:    http.StatusUnauthorized,
		errorType:     UnauthorizedErrorType,
	}
}
