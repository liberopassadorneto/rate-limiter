package middleware

import (
	"github.com/liberopassadorneto/rate-limiter/limiter"
	"net"
	"net/http"
	"strings"
)

func RateLimiterMiddleware(rl *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Invalid IP address", http.StatusInternalServerError)
				return
			}

			token := ""
			apiKey := r.Header.Get("API_KEY")
			if apiKey != "" {
				token = strings.TrimSpace(apiKey)
			}

			allowed, _, err := rl.Allow(ip, token)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
