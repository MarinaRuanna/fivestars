package customerror

import "net/http"

// StatusCodeFromError returns the HTTP status code for known CustomError, or 500.
func StatusCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var ce *CustomError
	if As(err, &ce) {
		return ce.StatusCode()
	}
	return http.StatusInternalServerError
}

// As checks if err is a *CustomError and assigns it to target.
func As(err error, target **CustomError) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*CustomError); ok {
		*target = e
		return true
	}
	return false
}
