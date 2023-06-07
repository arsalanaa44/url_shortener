package main

import (
	"github.com/mehrdad3301/url-shortner/internal/config"
	"github.com/mehrdad3301/url-shortner/internal/server"
)

func main() {
	cfg := config.New()
	server.Start(cfg)
}
