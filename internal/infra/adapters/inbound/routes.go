package inbound

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"fivestars/internal/infra/adapters/inbound/controller"
)

// Handlers is a collection of HTTP handlers used to wire routes.
type Handlers struct {
	Health         http.Handler
	Auth           *controller.AuthHandler
	User           *controller.UserHandler
	Establishments *controller.EstablishmentsHandler
}

// CreateChiRoutes registers routes using the chi router and returns it.
func CreateChiRoutes(h Handlers) http.Handler {
	r := chi.NewRouter()

	// Apply CORS first so OPTIONS/preflight requests are handled before
	// header validation (preflight doesn't include most headers).
	r.Use(CORS)

	// Register health endpoint (no extra header requirements)
	if h.Health != nil {
		r.Handle("/health", h.Health)
	}

	// Auth: /auth/register and /auth/login
	if h.Auth != nil {
		r.Post("/auth/register", h.Auth.Register)
		r.Post("/auth/login", h.Auth.Login)
	}

	// Users: /users/me requires Authorization header
	if h.User != nil {
		r.With(HeaderValidator(map[string]string{"Authorization": ""})).Get("/users/me", h.User.Me)
	}

	// Establishments: list endpoint
	if h.Establishments != nil {
		r.Get("/establishments", h.Establishments.List)
	}

	return r
}

// HeaderValidator returns a chi middleware that enforces required headers and
// optional expected values. Use an empty expected value to require presence
// only (e.g. map[string]string{"Authorization": ""}). For Content-Type
// exact match, set map["Content-Type"] = "application/json".
func HeaderValidator(required map[string]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// OPTIONS should be handled by CORS and not be rejected here
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			for k, expect := range required {
				val := r.Header.Get(k)
				if val == "" {
					http.Error(w, "missing header: "+k, http.StatusBadRequest)
					return
				}
				if expect != "" && val != expect {
					http.Error(w, "invalid header value: "+k, http.StatusBadRequest)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
