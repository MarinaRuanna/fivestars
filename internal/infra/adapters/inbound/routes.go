package inbound

import (
	"net/http"
	"time"

	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/adapters/inbound/controller"
	"fivestars/internal/infra/auth"

	"github.com/go-chi/chi/v5"
)

// Handlers is a collection of HTTP handlers used to wire routes.
type Handlers struct {
	Health         http.Handler
	Auth           *controller.AuthHandler
	User           *controller.UserHandler
	Establishments *controller.EstablishmentsHandler
	Checkins       *controller.CheckinsHandler
}

// CreateChiRoutes registers routes using the chi router and returns it.
func CreateChiRoutes(h Handlers, jwtSecret string, corsAllowedOrigins []string) http.Handler {
	r := chi.NewRouter()

	// Recover from panics and keep JSON error contract.
	r.Use(RecoverPanic)

	// Apply CORS first so OPTIONS/preflight requests are handled before
	// header validation (preflight doesn't include most headers).
	r.Use(CORS(corsAllowedOrigins))

	// Register health endpoint (no extra header requirements)
	if h.Health != nil {
		r.Handle("/health", h.Health)
	}

	// Auth: /auth/register and /auth/login
	if h.Auth != nil {
		r.Post("/auth/register", WithErrorEncoder(h.Auth.Register))
		r.With(LoginRateLimit(30, 10, 5*time.Minute)).Post("/auth/login", WithErrorEncoder(h.Auth.Login))
	}

	// Users: /users/me requires a valid bearer token
	if h.User != nil {
		r.With(auth.RequireAuth(jwtSecret)).Get("/users/me", WithErrorEncoder(h.User.GetUser))
	}

	// Establishments: list endpoint
	if h.Establishments != nil {
		r.Get("/establishments", WithErrorEncoder(h.Establishments.ListEstablishments))
	}

	// Checkins: create (protected) and list user's checkins
	if h.Checkins != nil {
		r.With(auth.RequireAuth(jwtSecret)).Post("/checkins", WithErrorEncoder(h.Checkins.CreateCheckin))
		r.With(auth.RequireAuth(jwtSecret)).Get("/checkins/me", WithErrorEncoder(h.Checkins.ListMyCheckins))
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
					controller.EncodeError(w, customerror.NewValidationError("missing header: "+k))
					return
				}
				if expect != "" && val != expect {
					controller.EncodeError(w, customerror.NewValidationError("invalid header value: "+k))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
