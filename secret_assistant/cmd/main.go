package main

import (
	"flag"
	"fmt"
	"log"
	"net"

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
	log.Printf("Received: %v", in.GetBody())
	return &api.ResponseMessage{Body: in.GetBody()}, nil
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
