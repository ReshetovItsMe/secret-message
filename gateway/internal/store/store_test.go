package store

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/cryptoclient"
)

func newTestStore(t *testing.T, ttl time.Duration) (Store, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	log := slog.New(slog.DiscardHandler)
	return New(client, ttl, log), mr
}

func testPayload() *cryptoclient.EncryptedPayload {
	return &cryptoclient.EncryptedPayload{
		PrivateKey:   []byte("private-key-bytes"),
		EncryptedKey: []byte("encrypted-key-bytes"),
		Data:         []byte("ciphertext"),
	}
}

func TestSendAndConsumeRoundTrip(t *testing.T) {
	ctx := context.Background()
	s, _ := newTestStore(t, time.Hour)

	id, err := s.Send(ctx, testPayload())
	require.NoError(t, err)
	require.NotEmpty(t, id)

	got, err := s.Consume(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, testPayload(), got)
}

func TestConsumeTwiceReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	s, _ := newTestStore(t, time.Hour)

	id, err := s.Send(ctx, testPayload())
	require.NoError(t, err)

	_, err = s.Consume(ctx, id)
	require.NoError(t, err)

	_, err = s.Consume(ctx, id)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestConsumeMissingReturnsNotFound(t *testing.T) {
	ctx := context.Background()
	s, _ := newTestStore(t, time.Hour)

	_, err := s.Consume(ctx, "no-such-id")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestTTLIsApplied(t *testing.T) {
	ctx := context.Background()
	s, mr := newTestStore(t, 42*time.Second)

	id, err := s.Send(ctx, testPayload())
	require.NoError(t, err)

	ttl := mr.TTL("msg:" + id)
	assert.Equal(t, 42*time.Second, ttl)
}

func TestKeysAreNamespaced(t *testing.T) {
	ctx := context.Background()
	s, _ := newTestStore(t, time.Hour)

	id, err := s.Send(ctx, testPayload())
	require.NoError(t, err)

	keys := make([]string, 0)
	iter := s.(*redisStore).client.Scan(ctx, 0, "*", 100).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	require.NoError(t, iter.Err())
	assert.Equal(t, []string{"msg:" + id}, keys)
}
