package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	DefaultNodePort         = 2222
	DefaultInternalRestPort = 61001
	DefaultLogLevel         = "info"
)

var (
	ErrConfigSecretKeyRequired = errors.New("SECRET_KEY environment variable is required")
)

type Config struct {
	SecretKey        string
	NodePort         int
	InternalRestPort int
	LogLevel         string

	Payload *NodePayload
}

func Load() (*Config, error) {
	cfg := &Config{
		NodePort:         DefaultNodePort,
		InternalRestPort: DefaultInternalRestPort,
		LogLevel:         DefaultLogLevel,
	}

	loadFromEnv(cfg)

	if cfg.SecretKey == "" {
		return nil, ErrConfigSecretKeyRequired
	}

	payload, err := ParseSecretKey(cfg.SecretKey)
	if err != nil {
		return nil, err
	}
	cfg.Payload = payload

	return cfg, nil
}

func loadFromEnv(cfg *Config) {
	if v := os.Getenv("SECRET_KEY"); v != "" {
		cfg.SecretKey = v
	}
	if v := os.Getenv("NODE_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil && port > 0 {
			cfg.NodePort = port
		}
	}
	if v := os.Getenv("INTERNAL_REST_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil && port > 0 {
			cfg.InternalRestPort = port
		}
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
}
