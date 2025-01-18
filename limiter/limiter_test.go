package limiter

import (
	"github.com/liberopassadorneto/rate-limiter/config"
	"github.com/liberopassadorneto/rate-limiter/strategy"
	"testing"
	"time"

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
