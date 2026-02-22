package auth

import (
	"net/http"
	"strings"
)

// RequireAuth retorna um middleware que exige Bearer JWT válido e coloca o user_id no context.
// Se não houver token ou for inválido, responde 401 e não chama o próximo handler.
func RequireAuth(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			const prefix = "Bearer "
			if !strings.HasPrefix(header, prefix) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"missing or invalid authorization header"}`))
				return
			}
			tokenString := strings.TrimPrefix(header, prefix)
			userID, err := ParseToken(tokenString, secret)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"invalid or expired token"}`))
				return
			}
			ctx := WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
