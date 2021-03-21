package main

import (
	"context"
	"log"
	"net"

	pb "bphun/k8sAutoscaling/TransactionAPI/TransactionAPI"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedTransactionAPIServer
}

func (s *server) SaveTransaction(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionReply, error) {
	log.Printf("Received: %v, Result: %v", in.GetInArr(), in.GetOutArr())
	return &pb.TransactionReply{Message: "Hello"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening on port %s", port)

	s := grpc.NewServer()
	pb.RegisterTransactionAPIServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
