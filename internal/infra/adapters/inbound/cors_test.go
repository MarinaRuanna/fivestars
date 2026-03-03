package inbound

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS_AllowsConfiguredOrigin(t *testing.T) {
	mw := CORS([]string{"http://localhost:3000"})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	mw(next).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "http://localhost:3000", rec.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_BlocksPreflightForDisallowedOrigin(t *testing.T) {
	mw := CORS([]string{"http://localhost:3000"})

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodOptions, "/auth/login", nil)
	req.Header.Set("Origin", "http://evil.example")
	rec := httptest.NewRecorder()

	mw(next).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}
