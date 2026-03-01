package customerror

import (
	"errors"
	"fmt"
)

type CustomError struct {
	messagePrefix string
	message       string
	statusCode    int
	errorType     ErrorType
}

type ErrorType string

func (e CustomError) Error() string {
	return fmt.Sprintf("%s: %s", e.messagePrefix, e.message)
}

func (e CustomError) ErrorType() ErrorType {
	return e.errorType
}

func (e CustomError) StatusCode() int {
	return e.statusCode
}

// TypeOf returns the domain error type when err (or wrapped err) is a CustomError.
func TypeOf(err error) (ErrorType, bool) {
	if err == nil {
		return "", false
	}

	var ce *CustomError
	if !errors.As(err, &ce) {
		return "", false
	}

	return ce.ErrorType(), true
}

// StatusCodeOf returns the status code when err (or wrapped err) is a CustomError.
func StatusCodeOf(err error) (int, bool) {
	if err == nil {
		return 0, false
	}

	var ce *CustomError
	if !errors.As(err, &ce) {
		return 0, false
	}

	return ce.StatusCode(), true
}
