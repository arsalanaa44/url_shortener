package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	expiration time.Duration
}

// NewRedisClient creates a new Redis client with the provided configuration.
func NewRedisClient(config Config) (*RedisClient, error) {
	// Create a new Redis client.
	client := redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: "",
		DB:       0,
	})

	// Test the connection to the Redis server.
	err := client.Ping(client.Context()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	// Create a new Redis client instance and return it.
	return &RedisClient{
		client:     client,
		expiration: config.Expiration * time.Minute,
	}, nil
}

// Set sets a key-value pair in Redis with the given expiration time.
func (c *RedisClient) Set(key string, value interface{}) error {
	err := c.client.Set(c.client.Context(), key, value, c.expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set %s: %v", key, err)
	}
	return nil
}

// Get retrieves a value from Redis by its key.
func (c *RedisClient) Get(key string) (string, error) {
	val, err := c.client.Get(c.client.Context(), key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %v", key, err)
	}
	return val, nil
}