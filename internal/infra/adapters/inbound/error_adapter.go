package inbound

import (
	"net/http"

	"fivestars/internal/infra/adapters/inbound/controller"
)

// HandlerWithError is an HTTP handler that can return an error to be encoded.
type HandlerWithError func(http.ResponseWriter, *http.Request) error

// WithErrorEncoder wraps a HandlerWithError and encodes returned errors as JSON.
func WithErrorEncoder(next HandlerWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			controller.EncodeError(w, err)
		}
	}
}
