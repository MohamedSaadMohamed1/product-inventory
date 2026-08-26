package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	sharedHttp "product-inventory/internal/http"
)

type client struct {
	limiter  *tokenBucket
	lastSeen time.Time
}

type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func newTokenBucket(maxTokens, refillRate float64) *tokenBucket {
	return &tokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *tokenBucket) allow() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens = tb.tokens + elapsed*tb.refillRate
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

// RateLimit creates an IP-based token bucket rate limiter middleware.
// limit: number of tokens generated per second.
// burst: maximum capacity of the bucket (maximum burst of requests).
func RateLimit(limit float64, burst float64) func(http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Background sweep goroutine to clean up idle records and prevent memory leak
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			mu.Lock()
			c, exists := clients[ip]
			if !exists {
				c = &client{limiter: newTokenBucket(burst, limit)}
				clients[ip] = c
			}
			c.lastSeen = time.Now()

			allowed := c.limiter.allow()
			mu.Unlock()

			if !allowed {
				sharedHttp.WriteError(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests. Please try again later.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
