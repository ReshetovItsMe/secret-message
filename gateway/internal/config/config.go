// Package config loads gateway configuration from environment variables.
//
// Environment variable names match the legacy Node.js gateway so that
// docker-compose and existing deployments keep working unchanged.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration for the gateway.
type Config struct {
	Host string
	Port int

	RedisHost     string
	RedisPort     int
	RedisUsername string
	RedisPassword string
	RedisDB       int

	// SecretAssistantURL is the host:port of the gRPC crypto service.
	SecretAssistantURL string

	// MessageTTL is how long a stored message lives before Redis evicts it.
	MessageTTL time.Duration

	// LogPretty switches the logger to a human-readable text format.
	LogPretty bool
}

// Load reads configuration from the environment, applying defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Host:               envStr("HOST", "0.0.0.0"),
		RedisHost:          envStr("REDIS_HOST", "redis"),
		RedisUsername:      envStr("REDIS_USERNAME", ""),
		RedisPassword:      envStr("REDIS_PASSWORD", ""),
		SecretAssistantURL: envStr("SECRET_ASSISTANT_URL", "secret-assistant:50052"),
	}

	var err error
	if cfg.Port, err = envInt("PORT", 3000); err != nil {
		return nil, err
	}
	if cfg.RedisPort, err = envInt("REDIS_PORT", 6379); err != nil {
		return nil, err
	}
	if cfg.RedisDB, err = envInt("REDIS_DB", 0); err != nil {
		return nil, err
	}
	if cfg.MessageTTL, err = envDuration("MESSAGE_TTL", 24*time.Hour); err != nil {
		return nil, err
	}
	if cfg.LogPretty, err = envBool("LOG_PRETTY", false); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Addr returns the listen address in host:port form.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func envStr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func envInt(key string, def int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parse env %s: %w", key, err)
	}
	return n, nil
}

func envBool(key string, def bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("parse env %s: %w", key, err)
	}
	return b, nil
}

func envDuration(key string, def time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("parse env %s: %w", key, err)
	}
	return d, nil
}
