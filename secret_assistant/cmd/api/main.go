package main

import (
	"flag"
	"fmt"
	pb "github.com/ReshetovItsMe/one-time-messaging-exchange-be/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedSecretAssistantServer
}

func (s *server) EncryptMessage(ctx context.Context, in *pb.RequestMessage) (*pb.ResponseMessage, error) {
	log.Printf("Received: %v", in.GetBody())
	return &pb.ResponseMessage{Id: "test", EncryptedMessage: in.GetBody()}, nil
}

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSecretAssistantServer(grpcServer, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
