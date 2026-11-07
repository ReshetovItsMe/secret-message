package main

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	dec "github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/decrypt"
	enc "github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/encrypt"
	"github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/keys"
	api "github.com/ReshetovItsMe/one-time-messaging-exchange-be/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

type server struct {
	api.UnimplementedSecretAssistantServer
	privateKey *rsa.PrivateKey
}

func (s *server) Encrypt(ctx context.Context, in *api.EncryptRequestMessage) (*api.EncryptedMessageResponse, error) {
	blob, err := enc.Encrypt(&s.privateKey.PublicKey, in.GetBody())
	if err != nil {
		return nil, err
	}
	return &api.EncryptedMessageResponse{Body: blob}, nil
}

func (s *server) Decrypt(ctx context.Context, in *api.DecryptRequestMessage) (*api.DecryptedMessageResponse, error) {
	plaintext, err := dec.Decrypt(s.privateKey, in.GetBody())
	if err != nil {
		return nil, err
	}
	return &api.DecryptedMessageResponse{Body: plaintext}, nil
}

func main() {
	flag.Parse()

	keyPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	if keyPath == "" {
		log.Fatalf("RSA_PRIVATE_KEY_PATH is not set: supply the path to the service RSA private key PEM file")
	}
	privateKey, err := keys.Load(keyPath)
	if err != nil {
		log.Fatalf("failed to load service private key: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serverInstance := server{privateKey: privateKey}
	grpcServer := grpc.NewServer()
	api.RegisterSecretAssistantServer(grpcServer, &serverInstance)
	reflection.Register(grpcServer)
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
