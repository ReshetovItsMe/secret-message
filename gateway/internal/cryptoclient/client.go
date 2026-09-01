// Package cryptoclient wraps the secret-assistant gRPC service.
//
// It is the gateway's adapter at the crypto seam: the rest of the codebase
// deals with EncryptedPayload, not protobuf messages.
package cryptoclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	messagepb "github.com/ReshetovItsMe/secret-message/gateway/internal/messagepb"
)

// EncryptedPayload is the wire shape shared with the crypto service.
// JSON field names match the secret-assistant contract (privateKey,
// encryptedKey, data).
type EncryptedPayload struct {
	PrivateKey   []byte `json:"privateKey"`
	EncryptedKey []byte `json:"encryptedKey"`
	Data         []byte `json:"data"`
}

// Client is the interface to the crypto service.
type Client interface {
	// Encrypt encrypts plaintext and returns the sealed payload.
	Encrypt(ctx context.Context, plaintext string) (*EncryptedPayload, error)
	// Decrypt reverses Encrypt.
	Decrypt(ctx context.Context, payload *EncryptedPayload) (string, error)
	// Close releases the underlying connection.
	Close() error
}

type grpcClient struct {
	conn *grpc.ClientConn
	svc  messagepb.SecretAssistantClient
}

// New dials the crypto service at addr (host:port). The connection is lazy:
// no I/O happens until the first RPC.
func New(addr string) (Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // in-cluster network only
	)
	if err != nil {
		return nil, fmt.Errorf("dial secret-assistant %s: %w", addr, err)
	}
	return &grpcClient{conn: conn, svc: messagepb.NewSecretAssistantClient(conn)}, nil
}

func (c *grpcClient) Encrypt(ctx context.Context, plaintext string) (*EncryptedPayload, error) {
	resp, err := c.svc.Encrypt(ctx, &messagepb.EncryptRequestMessage{Body: plaintext})
	if err != nil {
		return nil, fmt.Errorf("rpc encrypt: %w", err)
	}

	// Unmarshal of the response body cannot fail in practice: the shape is
	// fixed by the contract, so lo.Must keeps the happy path readable.
	var payload EncryptedPayload
	lo.Must0(json.Unmarshal(resp.GetBody(), &payload))
	return &payload, nil
}

func (c *grpcClient) Decrypt(ctx context.Context, payload *EncryptedPayload) (string, error) {
	body := lo.Must(json.Marshal(payload)) // struct of []byte fields — cannot fail

	resp, err := c.svc.Decrypt(ctx, &messagepb.DecryptRequestMessage{Body: body})
	if err != nil {
		return "", fmt.Errorf("rpc decrypt: %w", err)
	}
	return resp.GetBody(), nil
}

func (c *grpcClient) Close() error {
	return c.conn.Close()
}
