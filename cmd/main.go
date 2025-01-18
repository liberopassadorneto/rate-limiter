package main

import (
	"fmt"
	"github.com/liberopassadorneto/rate-limiter/config"
	"github.com/liberopassadorneto/rate-limiter/limiter"
	"github.com/liberopassadorneto/rate-limiter/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	redisStrategy := limiter.NewRedisLimiter(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)

	rateLimiter := limiter.NewRateLimiter(cfg, redisStrategy)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	handler := middleware.RateLimiterMiddleware(rateLimiter)(mux)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
