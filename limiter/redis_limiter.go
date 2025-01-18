package limiter

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLimiter struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisLimiter(address, password string, db int) *RedisLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return &RedisLimiter{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisLimiter) Increment(key string, window int) (int, error) {
	count, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		r.client.Expire(r.ctx, key, time.Duration(window)*time.Second)
	}
	return int(count), nil
}

func (r *RedisLimiter) Block(key string, duration int) error {
	return r.client.Set(r.ctx, key, "blocked", time.Duration(duration)*time.Second).Err()
}

func (r *RedisLimiter) IsBlocked(key string) (bool, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "blocked", nil
}
