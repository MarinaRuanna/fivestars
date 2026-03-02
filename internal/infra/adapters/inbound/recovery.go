package inbound

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"fivestars/internal/infra/adapters/inbound/httperror"
)

// RecoverPanic converts unexpected panics into a standardized 500 response.
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered on %s %s: %v\n%s", r.Method, r.URL.Path, rec, debug.Stack())
				httperror.Encode(w, fmt.Errorf("panic recovered: %v", rec))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
