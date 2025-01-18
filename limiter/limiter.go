package limiter

import (
	"github.com/liberopassadorneto/rate-limiter/config"
	"github.com/liberopassadorneto/rate-limiter/strategy"
)

type RateLimiter struct {
	config   *config.Config
	strategy strategy.LimiterStrategy
}

func NewRateLimiter(cfg *config.Config, strat strategy.LimiterStrategy) *RateLimiter {
	return &RateLimiter{
		config:   cfg,
		strategy: strat,
	}
}

func (rl *RateLimiter) Allow(ip, token string) (bool, string, error) {
	if token != "" {
		blocked, err := rl.strategy.IsBlocked(token)
		if err != nil {
			return false, "", err
		}
		if blocked {
			return false, "token", nil
		}

		count, err := rl.strategy.Increment(token, int(rl.config.TokenRateLimitWindow.Seconds()))
		if err != nil {
			return false, "", err
		}
		if count > rl.config.TokenRateLimit {
			rl.strategy.Block(token, int(rl.config.TokenBlockDuration.Seconds()))
			return false, "token", nil
		}
		return true, "", nil
	}

	blocked, err := rl.strategy.IsBlocked(ip)
	if err != nil {
		return false, "", err
	}
	if blocked {
		return false, "ip", nil
	}

	count, err := rl.strategy.Increment(ip, int(rl.config.IPRateLimitWindow.Seconds()))
	if err != nil {
		return false, "", err
	}
	if count > rl.config.IPRateLimit {
		err := rl.strategy.Block(ip, int(rl.config.IPBlockDuration.Seconds()))
		if err != nil {
			return false, "", err
		}
		return false, "ip", nil
	}
	return true, "", nil
}
