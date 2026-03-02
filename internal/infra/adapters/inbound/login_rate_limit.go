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
	mu           sync.Mutex
	ipCounters   map[string]rateLimitCounter
	emailCounter map[string]rateLimitCounter
	ipLimit      int
	emailLimit   int
	window       time.Duration
}

func newLoginRateLimiter(ipLimit, emailLimit int, window time.Duration) *loginRateLimiter {
	return &loginRateLimiter{
		ipCounters:   make(map[string]rateLimitCounter),
		emailCounter: make(map[string]rateLimitCounter),
		ipLimit:      ipLimit,
		emailLimit:   emailLimit,
		window:       window,
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

	if email != "" {
		emailCounter := l.emailCounter[email]
		if now.After(emailCounter.ResetAt) {
			emailCounter = rateLimitCounter{ResetAt: now.Add(l.window)}
		}
		emailCounter.Count++
		l.emailCounter[email] = emailCounter
		if emailCounter.Count > l.emailLimit {
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
	for k, v := range l.emailCounter {
		if now.After(v.ResetAt) {
			delete(l.emailCounter, k)
		}
	}
}

// LoginRateLimit applies in-memory rate limiting for login attempts by both IP and email.
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
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return xrip
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
}
