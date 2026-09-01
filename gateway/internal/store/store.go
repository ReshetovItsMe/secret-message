// Package store implements the message storage module.
//
// It is a deep module: callers interact through two methods (Send/Consume)
// and know nothing about the underlying key layout. Messages live under a
// single Redis key with a TTL and are atomically deleted on read, which is
// what provides the one-time guarantee advertised by the product.
package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/cryptoclient"
)

// ErrNotFound is returned by Consume when the message does not exist or has
// already been consumed (deleted).
var ErrNotFound = errors.New("message not found")

const keyPrefix = "msg:"

// Store is the storage interface for encrypted messages.
type Store interface {
	// Send stores payload and returns a fresh message id.
	Send(ctx context.Context, payload *cryptoclient.EncryptedPayload) (string, error)
	// Consume atomically reads and deletes the message with the given id.
	// Returns ErrNotFound if the message is missing or already consumed.
	Consume(ctx context.Context, id string) (*cryptoclient.EncryptedPayload, error)
	// Close releases the underlying connection.
	Close() error
}

type redisStore struct {
	client *redis.Client
	ttl    time.Duration
	log    *slog.Logger
}

// New creates a Store backed by Redis.
func New(client *redis.Client, ttl time.Duration, log *slog.Logger) Store {
	return &redisStore{client: client, ttl: ttl, log: log}
}

func (s *redisStore) Send(ctx context.Context, payload *cryptoclient.EncryptedPayload) (string, error) {
	id := uuid.NewString()
	key := keyPrefix + id

	blob, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	if err := s.client.Set(ctx, key, blob, s.ttl).Err(); err != nil {
		return "", fmt.Errorf("redis set: %w", err)
	}

	s.log.Debug("message stored", "id", id, "ttl", s.ttl)
	return id, nil
}

func (s *redisStore) Consume(ctx context.Context, id string) (*cryptoclient.EncryptedPayload, error) {
	key := keyPrefix + id

	// GETDEL is atomic: read + delete in one round trip, so two concurrent
	// consumers can never both see the message (the one-time guarantee).
	blob, err := s.client.GetDel(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("redis getdel: %w", err)
	}

	var payload cryptoclient.EncryptedPayload
	if err := json.Unmarshal(blob, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal payload: %w", err)
	}

	s.log.Debug("message consumed", "id", id)
	return &payload, nil
}

func (s *redisStore) Close() error {
	return s.client.Close()
}
