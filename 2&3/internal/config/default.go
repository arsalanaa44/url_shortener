package config

import (
	"time"

	"github.com/mehrdad3301/url-shortner/internal/redis"
)

func Default() Config { 
	return Config{ 
		Port: 8080, 
		Database: redis.Config{
			URL: "redis:6379", 
			Expiration: time.Minute * 5,
		},
	} 
}