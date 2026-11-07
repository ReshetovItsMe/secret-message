package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/store"
)

type fakeCrypto struct {
	blob        []byte
	encryptErr  error
	decryptText string
	decryptErr  error
}

func (f *fakeCrypto) Encrypt(_ context.Context, _ string) ([]byte, error) {
	return f.blob, f.encryptErr
}

func (f *fakeCrypto) Decrypt(_ context.Context, _ []byte) (string, error) {
	return f.decryptText, f.decryptErr
}

func (f *fakeCrypto) Close() error { return nil }

type fakeStore struct {
	id         string
	sendErr    error
	blob       []byte
	consumeErr error
}

func (f *fakeStore) Send(_ context.Context, _ []byte) (string, error) {
	return f.id, f.sendErr
}

func (f *fakeStore) Consume(_ context.Context, _ string) ([]byte, error) {
	if f.consumeErr != nil {
		return nil, f.consumeErr
	}
	return f.blob, nil
}

func (f *fakeStore) Close() error { return nil }

func TestCreateMessage(t *testing.T) {
	ctx := context.Background()
	crypto := &fakeCrypto{blob: []byte("x")}
	storage := &fakeStore{id: "msg-1"}
	svc := New(crypto, storage)

	id, err := svc.CreateMessage(ctx, "hello")
	require.NoError(t, err)
	assert.Equal(t, "msg-1", id)
}

func TestCreateMessagePropagatesEncryptError(t *testing.T) {
	ctx := context.Background()
	crypto := &fakeCrypto{encryptErr: errors.New("crypto down")}
	svc := New(crypto, &fakeStore{})

	_, err := svc.CreateMessage(ctx, "hello")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "encrypt message")
}

func TestCreateMessagePropagatesStoreError(t *testing.T) {
	ctx := context.Background()
	crypto := &fakeCrypto{blob: []byte{}}
	storage := &fakeStore{sendErr: errors.New("redis down")}
	svc := New(crypto, storage)

	_, err := svc.CreateMessage(ctx, "hello")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "store message")
}

func TestGetMessage(t *testing.T) {
	ctx := context.Background()
	crypto := &fakeCrypto{decryptText: "secret text"}
	storage := &fakeStore{blob: []byte("c")}
	svc := New(crypto, storage)

	text, err := svc.GetMessage(ctx, "msg-1")
	require.NoError(t, err)
	assert.Equal(t, "secret text", text)
}

func TestGetMessageMapsNotFound(t *testing.T) {
	ctx := context.Background()
	storage := &fakeStore{consumeErr: store.ErrNotFound}
	svc := New(&fakeCrypto{}, storage)

	_, err := svc.GetMessage(ctx, "missing")
	assert.ErrorIs(t, err, ErrMessageNotFound)
}

func TestGetMessagePropagatesDecryptError(t *testing.T) {
	ctx := context.Background()
	crypto := &fakeCrypto{decryptErr: errors.New("bad key")}
	storage := &fakeStore{blob: []byte{}}
	svc := New(crypto, storage)

	_, err := svc.GetMessage(ctx, "msg-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decrypt message")
}
