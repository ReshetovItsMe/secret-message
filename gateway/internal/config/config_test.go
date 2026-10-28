package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOverrides(t *testing.T) {
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_HOST", "myredis")
	t.Setenv("REDIS_PORT", "6380")
	t.Setenv("REDIS_DB", "3")
	t.Setenv("SECRET_ASSISTANT_URL", "crypto:6000")
	t.Setenv("MESSAGE_TTL", "1h30m")
	t.Setenv("LOG_PRETTY", "true")

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "myredis", cfg.RedisHost)
	assert.Equal(t, 6380, cfg.RedisPort)
	assert.Equal(t, 3, cfg.RedisDB)
	assert.Equal(t, "crypto:6000", cfg.SecretAssistantURL)
	assert.Equal(t, 90*time.Minute, cfg.MessageTTL)
	assert.True(t, cfg.LogPretty)
	assert.Equal(t, "127.0.0.1:8080", cfg.Addr())
}

func TestLoadInvalidPort(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "PORT")
}

func TestLoadInvalidTTL(t *testing.T) {
	t.Setenv("MESSAGE_TTL", "soon")

	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "MESSAGE_TTL")
}

func TestLoadInvalidBool(t *testing.T) {
	t.Setenv("LOG_PRETTY", "yes")

	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "LOG_PRETTY")
}
