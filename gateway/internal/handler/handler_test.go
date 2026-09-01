package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/cryptoclient"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/handler"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/server"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/service"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/store"
)

type stubCrypto struct {
	payload *cryptoclient.EncryptedPayload
	text    string
}

func (s *stubCrypto) Encrypt(_ context.Context, _ string) (*cryptoclient.EncryptedPayload, error) {
	return s.payload, nil
}

func (s *stubCrypto) Decrypt(_ context.Context, _ *cryptoclient.EncryptedPayload) (string, error) {
	return s.text, nil
}

func (s *stubCrypto) Close() error { return nil }

type stubStore struct {
	payload *cryptoclient.EncryptedPayload
	id      string
	err     error
}

func (s *stubStore) Send(_ context.Context, _ *cryptoclient.EncryptedPayload) (string, error) {
	return s.id, s.err
}

func (s *stubStore) Consume(_ context.Context, _ string) (*cryptoclient.EncryptedPayload, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.payload, nil
}

func (s *stubStore) Close() error { return nil }

func newTestServer(t *testing.T, crypto cryptoclient.Client, storage store.Store) http.Handler {
	t.Helper()
	log := slog.New(slog.DiscardHandler)
	svc := service.New(crypto, storage)
	return server.New(handler.New(svc, log), log)
}

func doRequest(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreateMessageReturns201(t *testing.T) {
	h := newTestServer(t,
		&stubCrypto{payload: &cryptoclient.EncryptedPayload{Data: []byte("c")}},
		&stubStore{id: "abc-123"},
	)

	rec := doRequest(t, h, http.MethodPost, "/message", `{"message":"hello world"}`)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp struct {
		MessageID string `json:"messageId"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "abc-123", resp.MessageID)
}

func TestCreateMessageRejectsMissingField(t *testing.T) {
	h := newTestServer(t, &stubCrypto{}, &stubStore{})

	rec := doRequest(t, h, http.MethodPost, "/message", `{}`)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateMessageRejectsInvalidJSON(t *testing.T) {
	h := newTestServer(t, &stubCrypto{}, &stubStore{})

	rec := doRequest(t, h, http.MethodPost, "/message", `{not json`)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMessageReturns200(t *testing.T) {
	h := newTestServer(t,
		&stubCrypto{text: "secret"},
		&stubStore{payload: &cryptoclient.EncryptedPayload{Data: []byte("c")}},
	)

	rec := doRequest(t, h, http.MethodGet, "/message?messageId=abc", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "secret", resp.Message)
}

func TestGetMessageMissingReturns404(t *testing.T) {
	h := newTestServer(t, &stubCrypto{}, &stubStore{err: store.ErrNotFound})

	rec := doRequest(t, h, http.MethodGet, "/message?messageId=nope", "")
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	errObj, ok := body["error"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "NOT_FOUND", errObj["code"])
}

func TestGetMessageRequiresMessageID(t *testing.T) {
	h := newTestServer(t, &stubCrypto{}, &stubStore{})

	rec := doRequest(t, h, http.MethodGet, "/message", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUnknownRouteReturns404(t *testing.T) {
	h := newTestServer(t, &stubCrypto{}, &stubStore{})

	rec := doRequest(t, h, http.MethodGet, "/nope", "")
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestInternalErrorIsObfuscated(t *testing.T) {
	h := newTestServer(t,
		&stubCrypto{},
		&stubStore{err: errors.New("redis exploded with details")},
	)

	rec := doRequest(t, h, http.MethodPost, "/message", `{"message":"x"}`)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotContains(t, rec.Body.String(), "redis exploded")
}
