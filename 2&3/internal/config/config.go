package config

import (
	"github.com/mehrdad3301/url-shortner/internal/redis"
)

type Config struct {
	Port       int          `koanf:"port"`
	Database   redis.Config `koanf:"database"`
}
