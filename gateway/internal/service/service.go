// Package service contains the application's business logic: the
// encrypt → store → decrypt orchestration behind a small interface.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/cryptoclient"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/store"
)

// ErrMessageNotFound is returned by GetMessage when the message does not
// exist or has already been consumed.
var ErrMessageNotFound = errors.New("message not found")

// Service orchestrates the message lifecycle.
type Service struct {
	crypto cryptoclient.Client
	store  store.Store
}

// New wires the service with its adapters.
func New(crypto cryptoclient.Client, storage store.Store) *Service {
	return &Service{crypto: crypto, store: storage}
}

// CreateMessage encrypts plaintext, stores the sealed payload and returns a
// one-time message id.
func (s *Service) CreateMessage(ctx context.Context, plaintext string) (string, error) {
	payload, err := s.crypto.Encrypt(ctx, plaintext)
	if err != nil {
		return "", fmt.Errorf("encrypt message: %w", err)
	}

	id, err := s.store.Send(ctx, payload)
	if err != nil {
		return "", fmt.Errorf("store message: %w", err)
	}
	return id, nil
}

// GetMessage reads and deletes the message with the given id, then decrypts
// it. Returns ErrMessageNotFound if the message is missing or already
// consumed.
func (s *Service) GetMessage(ctx context.Context, id string) (string, error) {
	payload, err := s.store.Consume(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		return "", ErrMessageNotFound
	}
	if err != nil {
		return "", fmt.Errorf("load message: %w", err)
	}

	plaintext, err := s.crypto.Decrypt(ctx, payload)
	if err != nil {
		return "", fmt.Errorf("decrypt message: %w", err)
	}
	return plaintext, nil
}
