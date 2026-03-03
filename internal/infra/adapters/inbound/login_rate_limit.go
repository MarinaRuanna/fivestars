package inbound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/adapters/inbound/httperror"
)

const loginRateLimitBodyMaxBytes int64 = 1 << 20 // 1MB

type rateLimitCounter struct {
	Count   int
	ResetAt time.Time
}

type loginRateLimiter struct {
	mu                 sync.Mutex
	ipCounters         map[string]rateLimitCounter
	emailPerIPCounters map[string]rateLimitCounter
	ipLimit            int
	emailPerIPLimit    int
	window             time.Duration
}

func newLoginRateLimiter(ipLimit, emailLimit int, window time.Duration) *loginRateLimiter {
	return &loginRateLimiter{
		ipCounters:         make(map[string]rateLimitCounter),
		emailPerIPCounters: make(map[string]rateLimitCounter),
		ipLimit:            ipLimit,
		emailPerIPLimit:    emailLimit,
		window:             window,
	}
}

func (l *loginRateLimiter) hitAndCheck(ip, email string) bool {
	now := time.Now().UTC()
	l.mu.Lock()
	defer l.mu.Unlock()

	// Opportunistic cleanup keeps memory bounded.
	l.cleanupExpired(now)

	if ip != "" {
		ipCounter := l.ipCounters[ip]
		if now.After(ipCounter.ResetAt) {
			ipCounter = rateLimitCounter{ResetAt: now.Add(l.window)}
		}
		ipCounter.Count++
		l.ipCounters[ip] = ipCounter
		if ipCounter.Count > l.ipLimit {
			return false
		}
	}

	// Email throttling is scoped per IP to avoid global account lockout.
	if ip != "" && email != "" {
		key := ip + "|" + email
		emailCounter := l.emailPerIPCounters[key]
		if now.After(emailCounter.ResetAt) {
			emailCounter = rateLimitCounter{ResetAt: now.Add(l.window)}
		}
		emailCounter.Count++
		l.emailPerIPCounters[key] = emailCounter
		if emailCounter.Count > l.emailPerIPLimit {
			return false
		}
	}

	return true
}

func (l *loginRateLimiter) cleanupExpired(now time.Time) {
	for k, v := range l.ipCounters {
		if now.After(v.ResetAt) {
			delete(l.ipCounters, k)
		}
	}
	for k, v := range l.emailPerIPCounters {
		if now.After(v.ResetAt) {
			delete(l.emailPerIPCounters, k)
		}
	}
}

// LoginRateLimit applies in-memory rate limiting for login attempts by IP and by IP+email.
func LoginRateLimit(ipLimit, emailLimit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := newLoginRateLimiter(ipLimit, emailLimit, window)

	type loginRequest struct {
		Email string `json:"email"`
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(io.LimitReader(r.Body, loginRateLimitBodyMaxBytes))
			if err != nil {
				httperror.Encode(w, customerror.NewValidationError("invalid body"))
				return
			}
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			var req loginRequest
			_ = json.Unmarshal(body, &req)

			ip := clientIP(r)
			email := strings.ToLower(strings.TrimSpace(req.Email))

			if ok := limiter.hitAndCheck(ip, email); !ok {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
				httperror.Encode(w, customerror.NewTooManyRequestsError("too many login attempts, try again later"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	// Use socket address only. Forwarded headers are ignored here to prevent spoofing.
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
