package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Limiter: 2 requests per second, maximum burst of 3
	limiter := RateLimit(2, 3)

	handler := limiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// 1. Initial burst of 3 requests should all pass
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 2. The 4th request immediately after should be blocked (429)
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// 3. A different client IP should be allowed (isolated client limits)
	reqDifferent := httptest.NewRequest("GET", "/", nil)
	reqDifferent.RemoteAddr = "192.168.1.2:1234"
	recDifferent := httptest.NewRecorder()

	handler.ServeHTTP(recDifferent, reqDifferent)
	assert.Equal(t, http.StatusOK, recDifferent.Code)
}
