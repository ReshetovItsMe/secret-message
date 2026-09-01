// Package config loads gateway configuration from environment variables.
//
// Environment variable names match the legacy Node.js gateway so that
// docker-compose and existing deployments keep working unchanged.
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds all runtime configuration for the gateway.
type Config struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port int    `env:"PORT" envDefault:"3000"`

	RedisHost     string `env:"REDIS_HOST" envDefault:"redis"`
	RedisPort     int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisUsername string `env:"REDIS_USERNAME" envDefault:""`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`

	// SecretAssistantURL is the host:port of the gRPC crypto service.
	SecretAssistantURL string `env:"SECRET_ASSISTANT_URL" envDefault:"secret-assistant:50052"`

	// MessageTTL is how long a stored message lives before Redis evicts it.
	MessageTTL time.Duration `env:"MESSAGE_TTL" envDefault:"24h"`

	// LogPretty switches the logger to a human-readable text format.
	LogPretty bool `env:"LOG_PRETTY" envDefault:"false"`
}

// Load reads configuration from the environment.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}
	return cfg, nil
}

// Addr returns the listen address in host:port form.
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
