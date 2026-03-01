package httperror

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/domain/customerror"
)

// Encode writes a standardized JSON error response.
func Encode(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if statusCode, ok := customerror.StatusCodeOf(err); ok {
		code = statusCode
	}

	message := err.Error()
	if code >= http.StatusInternalServerError {
		message = "internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
