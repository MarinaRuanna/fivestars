package inbound

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_loginRateLimiter_EmailLimitIsScopedPerIP(t *testing.T) {
	limiter := newLoginRateLimiter(100, 2, time.Minute)
	email := "target@example.com"

	assert.True(t, limiter.hitAndCheck("10.0.0.1", email))
	assert.True(t, limiter.hitAndCheck("10.0.0.1", email))
	assert.False(t, limiter.hitAndCheck("10.0.0.1", email))

	// Same email from a different IP must not be blocked by attempts from another IP.
	assert.True(t, limiter.hitAndCheck("10.0.0.2", email))
}

func Test_clientIP_IgnoresForwardedHeaders(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/login", nil)
	req.RemoteAddr = "10.1.2.3:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req.Header.Set("X-Real-IP", "203.0.113.11")

	ip := clientIP(req)
	assert.Equal(t, "10.1.2.3", ip)
}
