package customerror

import "fmt"

type CustomError struct {
	messagePrefix string
	message       string
	statusCode    int
	errorType     ErrorType
}

type ErrorType string

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.messagePrefix, e.message)
}

func (e *CustomError) StatusCode() int {
	return e.statusCode
}

func (e *CustomError) ErrorType() ErrorType {
	return e.errorType
}
