package main

import (
	"flag"
	"fmt"
	"log"
	"net"

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

func (s *server) Encrypt(ctx context.Context, in *api.RequestMessage) (*api.ResponseMessage, error) {
	text := in.GetBody()
	log.Printf("Received: %v", text)
	ciphertext, err := enc.Encrypt(text)
	if err != nil {
		return nil, err
	}
	return &api.ResponseMessage{Body: ciphertext}, nil
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
