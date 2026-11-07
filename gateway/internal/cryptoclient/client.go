// Package cryptoclient wraps the secret-assistant gRPC service.
//
// It is the gateway's adapter at the crypto seam: the rest of the codebase
// deals with opaque []byte blobs, not protobuf messages. The gateway has no
// knowledge of the blob's internal layout (keys or ciphertext structure).
package cryptoclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	messagepb "github.com/ReshetovItsMe/secret-message/gateway/internal/messagepb"
)

// Client is the interface to the crypto service.
type Client interface {
	// Encrypt encrypts plaintext and returns the opaque sealed blob.
	Encrypt(ctx context.Context, plaintext string) ([]byte, error)
	// Decrypt reverses Encrypt, turning an opaque blob back into plaintext.
	Decrypt(ctx context.Context, blob []byte) (string, error)
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

func (c *grpcClient) Encrypt(ctx context.Context, plaintext string) ([]byte, error) {
	resp, err := c.svc.Encrypt(ctx, &messagepb.EncryptRequest{Plaintext: plaintext})
	if err != nil {
		return nil, fmt.Errorf("rpc encrypt: %w", err)
	}
	return resp.GetBlob(), nil
}

func (c *grpcClient) Decrypt(ctx context.Context, blob []byte) (string, error) {
	resp, err := c.svc.Decrypt(ctx, &messagepb.DecryptRequest{Blob: blob})
	if err != nil {
		return "", fmt.Errorf("rpc decrypt: %w", err)
	}
	return resp.GetPlaintext(), nil
}

func (c *grpcClient) Close() error {
	return c.conn.Close()
}
