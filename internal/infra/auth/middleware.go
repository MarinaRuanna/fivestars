package auth

import (
	"net/http"
	"strings"

	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/adapters/inbound/httperror"
)

// RequireAuth retorna um middleware que exige Bearer JWT válido e coloca o user_id no context.
// Se não houver token ou for inválido, responde 401 e não chama o próximo handler.
func RequireAuth(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if !strings.HasPrefix(header, prefix) {
				httperror.Encode(w, customerror.NewUnauthorizedError("missing or invalid authorization header"))
				return
			}
			tokenString := strings.TrimPrefix(header, prefix)
			userID, err := ParseToken(tokenString, secret)
			if err != nil {
				httperror.Encode(w, customerror.NewUnauthorizedError("invalid or expired token"))
				return
			}
			ctx := WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
