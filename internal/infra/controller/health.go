package controller

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler handles /health (checks DB with SELECT 1).
type HealthHandler struct {
	pool *pgxpool.Pool
}

// NewHealthHandler returns a new HealthHandler.
func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

// ServeHTTP responds 200 if DB is reachable, 503 otherwise.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	if err := h.pool.Ping(ctx); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

// Ping is a simple func for chi-style handlers that accept context.
func (h *HealthHandler) Ping(ctx context.Context) (int, []byte) {
	if err := h.pool.Ping(ctx); err != nil {
		return http.StatusServiceUnavailable, []byte(`{"status":"unavailable"}`)
	}
	return http.StatusOK, []byte(`{"status":"ok"}`)
}
