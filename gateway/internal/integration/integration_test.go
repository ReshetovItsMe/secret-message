// Package integration exercises the full gateway stack over real TCP:
// HTTP server + Redis (miniredis) + gRPC client against an in-process fake
// crypto service.
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/ReshetovItsMe/secret-message/gateway/internal/cryptoclient"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/handler"
	messagepb "github.com/ReshetovItsMe/secret-message/gateway/internal/messagepb"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/server"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/service"
	"github.com/ReshetovItsMe/secret-message/gateway/internal/store"
)

// fakeCryptoServer is a minimal in-process SecretAssistant implementing the
// gRPC contract with a reversible transform (no real encryption needed to
// prove the wiring works).
type fakeCryptoServer struct {
	messagepb.UnimplementedSecretAssistantServer
}

func (fakeCryptoServer) Encrypt(_ context.Context, in *messagepb.EncryptRequestMessage) (*messagepb.EncryptedMessageResponse, error) {
	payload := cryptoclient.EncryptedPayload{
		PrivateKey:   []byte("pk"),
		EncryptedKey: []byte("ek"),
		Data:         []byte("cipher:" + in.GetBody()),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &messagepb.EncryptedMessageResponse{Body: body}, nil
}

func (fakeCryptoServer) Decrypt(_ context.Context, in *messagepb.DecryptRequestMessage) (*messagepb.DecryptedMessageResponse, error) {
	var payload cryptoclient.EncryptedPayload
	if err := json.Unmarshal(in.GetBody(), &payload); err != nil {
		return nil, err
	}
	return &messagepb.DecryptedMessageResponse{Body: string(payload.Data)}, nil
}

// startFakeCrypto spins up the fake gRPC service on a random port.
func startFakeCrypto(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv := grpc.NewServer()
	messagepb.RegisterSecretAssistantServer(srv, fakeCryptoServer{})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

func TestEndToEnd(t *testing.T) {
	log := slog.New(slog.DiscardHandler)

	// Real Redis protocol via miniredis on TCP.
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.StartAddr("127.0.0.1:0"))
	t.Cleanup(mr.Close)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	cryptoAddr := startFakeCrypto(t)
	crypto, err := cryptoclient.New(cryptoAddr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = crypto.Close() })

	storage := store.New(rdb, time.Hour, log)
	svc := service.New(crypto, storage)
	h := server.New(handler.New(svc, log), log)

	// Real HTTP server on a random port.
	ts := httptest.NewServer(h)
	defer ts.Close()

	// 1. Create a message → 201 + messageId.
	createResp, err := http.Post(ts.URL+"/message", "application/json",
		bytes.NewBufferString(`{"message":"hello from e2e"}`))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)
	var created struct {
		MessageID string `json:"messageId"`
	}
	require.NoError(t, json.NewDecoder(createResp.Body).Decode(&created))
	createResp.Body.Close()
	require.NotEmpty(t, created.MessageID)

	// 2. Read it once → 200 + decrypted text.
	getResp, err := http.Get(fmt.Sprintf("%s/message?messageId=%s", ts.URL, created.MessageID))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	var got struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(getResp.Body).Decode(&got))
	getResp.Body.Close()
	assert.Equal(t, "cipher:hello from e2e", got.Message)

	// 3. Read it again → 404 (one-time guarantee).
	againResp, err := http.Get(fmt.Sprintf("%s/message?messageId=%s", ts.URL, created.MessageID))
	require.NoError(t, err)
	againResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, againResp.StatusCode)

	// 4. Validation → 400.
	badResp, err := http.Post(ts.URL+"/message", "application/json", bytes.NewBufferString(`{}`))
	require.NoError(t, err)
	badResp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, badResp.StatusCode)

	// 5. Healthz.
	healthResp, err := http.Get(ts.URL + "/healthz")
	require.NoError(t, err)
	healthResp.Body.Close()
	assert.Equal(t, http.StatusOK, healthResp.StatusCode)
}
