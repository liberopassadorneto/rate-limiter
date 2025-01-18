package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	IPRateLimit          int
	IPRateLimitWindow    time.Duration
	IPBlockDuration      time.Duration
	TokenRateLimit       int
	TokenRateLimitWindow time.Duration
	TokenBlockDuration   time.Duration
	RedisAddress         string
	RedisPassword        string
	RedisDB              int
	ServerPort           string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading configuration from environment variables")
	}

	ipRateLimit, _ := strconv.Atoi(getEnv("IP_RATE_LIMIT", "10"))
	ipRateLimitWindow, _ := time.ParseDuration(getEnv("IP_RATE_LIMIT_WINDOW", "1s"))
	ipBlockDuration, _ := time.ParseDuration(getEnv("IP_BLOCK_DURATION", "5m"))

	tokenRateLimit, _ := strconv.Atoi(getEnv("TOKEN_RATE_LIMIT", "100"))
	tokenRateLimitWindow, _ := time.ParseDuration(getEnv("TOKEN_RATE_LIMIT_WINDOW", "1s"))
	tokenBlockDuration, _ := time.ParseDuration(getEnv("TOKEN_BLOCK_DURATION", "5m"))

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		IPRateLimit:          ipRateLimit,
		IPRateLimitWindow:    ipRateLimitWindow,
		IPBlockDuration:      ipBlockDuration,
		TokenRateLimit:       tokenRateLimit,
		TokenRateLimitWindow: tokenRateLimitWindow,
		TokenBlockDuration:   tokenBlockDuration,
		RedisAddress:         getEnv("REDIS_ADDRESS", "localhost:6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		RedisDB:              redisDB,
		ServerPort:           getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
