package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	dec "github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/decrypt"
	enc "github.com/ReshetovItsMe/one-time-messaging-exchange-be/internal/encrypt"
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
}

func (s *server) Encrypt(ctx context.Context, in *api.EncryptRequestMessage) (*api.EncryptedMessageResponse, error) {
	text := in.GetBody()
	log.Printf("Received some message to encrypt")
	ciphertext, err := enc.Encrypt(text)
	if err != nil {
		return nil, err
	}

	return &api.EncryptedMessageResponse{Body: ciphertext.Data}, nil
}

func (s *server) Decrypt(ctx context.Context, in *api.DecryptRequestMessage) (*api.DecryptedMessageResponse, error) {
	text := in.GetBody()
	log.Printf("Received some message to decrypt")
	unwrappedText, err := dec.Decrypt(text)
	if err != nil {
		return nil, err
	}
	return &api.DecryptedMessageResponse{Body: unwrappedText}, nil
}

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	serverInstance := server{}
	grpcServer := grpc.NewServer()
	api.RegisterSecretAssistantServer(grpcServer, &serverInstance)
	reflection.Register(grpcServer)
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
