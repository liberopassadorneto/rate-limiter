package limiter

import (
	"testing"
	"time"

	"github.com/liberopassadorneto/rate-limiter/config"
	"github.com/liberopassadorneto/rate-limiter/strategy"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupRedisLimiter(t *testing.T) (strategy.LimiterStrategy, func()) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when starting miniredis", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	limiter := NewRedisLimiter(s.Addr(), "", 0)

	teardown := func() {
		rdb.Close()
		s.Close()
	}

	return limiter, teardown
}

func TestRateLimiter(t *testing.T) {
	limiterStrategy, teardown := setupRedisLimiter(t)
	defer teardown()

	cfg := &config.Config{
		IPRateLimit:          5,
		IPRateLimitWindow:    1 * time.Second,
		IPBlockDuration:      5 * time.Minute,
		TokenRateLimit:       10,
		TokenRateLimitWindow: 1 * time.Second,
		TokenBlockDuration:   5 * time.Minute,
	}

	rl := NewRateLimiter(cfg, limiterStrategy)

	ip := "192.168.1.1"
	token := "abc123"

	for i := 1; i <= 5; i++ {
		allowed, _, err := rl.Allow(ip, "")
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, limiterType, err := rl.Allow(ip, "")
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, "ip", limiterType)

	for i := 1; i <= 10; i++ {
		allowed, _, err := rl.Allow(ip, token)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, limiterType, err = rl.Allow(ip, token)
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, "token", limiterType)
}

func TestIPAllowed(t *testing.T) {
	limiterStrategy, teardown := setupRedisLimiter(t)
	defer teardown()

	cfg := &config.Config{
		IPRateLimit:       5,
		IPRateLimitWindow: 1 * time.Second,
	}

	rl := NewRateLimiter(cfg, limiterStrategy)
	ip := "192.168.1.1"

	// Testando o limite permitido para IP
	for i := 1; i <= 5; i++ {
		allowed, _, err := rl.Allow(ip, "")
		assert.NoError(t, err)
		assert.True(t, allowed)
	}
}

func TestIPBlocked(t *testing.T) {
	limiterStrategy, teardown := setupRedisLimiter(t)
	defer teardown()

	cfg := &config.Config{
		IPRateLimit:       5,
		IPRateLimitWindow: 1 * time.Second,
		IPBlockDuration:   2 * time.Second,
	}

	rl := NewRateLimiter(cfg, limiterStrategy)
	ip := "192.168.1.1"

	// Excedendo o limite para bloquear o IP
	for i := 1; i <= 5; i++ {
		allowed, _, err := rl.Allow(ip, "")
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, limiterType, err := rl.Allow(ip, "")
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, "ip", limiterType)

	blocked, err := rl.IsBlocked(ip)
	assert.NoError(t, err)
	assert.True(t, blocked)
}

func TestTokenAllowed(t *testing.T) {
	limiterStrategy, teardown := setupRedisLimiter(t)
	defer teardown()

	cfg := &config.Config{
		TokenRateLimit:       10,
		TokenRateLimitWindow: 1 * time.Second,
	}

	rl := NewRateLimiter(cfg, limiterStrategy)
	token := "abc123"

	// Testando o limite permitido para Token
	for i := 1; i <= 10; i++ {
		allowed, _, err := rl.Allow("", token)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}
}

func TestTokenBlocked(t *testing.T) {
	limiterStrategy, teardown := setupRedisLimiter(t)
	defer teardown()

	cfg := &config.Config{
		TokenRateLimit:       10,
		TokenRateLimitWindow: 1 * time.Second,
		TokenBlockDuration:   2 * time.Second,
	}

	rl := NewRateLimiter(cfg, limiterStrategy)
	token := "abc123"

	// Excedendo o limite para bloquear o Token
	for i := 1; i <= 10; i++ {
		allowed, _, err := rl.Allow("", token)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, limiterType, err := rl.Allow("", token)
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, "token", limiterType)

	blocked, err := rl.IsBlocked(token)
	assert.NoError(t, err)
	assert.True(t, blocked)
}

func TestIPAndTokenBlocked(t *testing.T) {
	limiterStrategy, teardown := setupRedisLimiter(t)
	defer teardown()

	cfg := &config.Config{
		IPRateLimit:          5,
		IPRateLimitWindow:    1 * time.Second,
		IPBlockDuration:      2 * time.Second,
		TokenRateLimit:       10,
		TokenRateLimitWindow: 1 * time.Second,
		TokenBlockDuration:   2 * time.Second,
	}

	rl := NewRateLimiter(cfg, limiterStrategy)
	ip := "192.168.1.1"
	token := "abc123"

	// Bloqueando o IP
	for i := 1; i <= 5; i++ {
		allowed, _, err := rl.Allow(ip, "")
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, limiterType, err := rl.Allow(ip, "")
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, "ip", limiterType)

	blocked, err := rl.IsBlocked(ip)
	assert.NoError(t, err)
	assert.True(t, blocked)

	// Bloqueando o Token
	for i := 1; i <= 10; i++ {
		allowed, _, err := rl.Allow(ip, token)
		assert.NoError(t, err)
		assert.True(t, allowed)
	}

	allowed, limiterType, err = rl.Allow(ip, token)
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, "token", limiterType)

	blocked, err = rl.IsBlocked(token)
	assert.NoError(t, err)
	assert.True(t, blocked)
}
