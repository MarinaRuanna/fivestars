package auth

import "context"

type contextKey string

const userIDKey contextKey = "user_id"

// WithUserID retorna um context com o user_id definido.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext retorna o user_id do context ou "" se não estiver autenticado.
func UserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}
