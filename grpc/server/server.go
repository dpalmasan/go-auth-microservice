package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/go-auth-microservice/grpc/auth"
	"github.com/go-auth-microservice/models"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server port")
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) ValidateToken(ctx context.Context, in *pb.TokenRequest) (*pb.TokenResponse, error) {
	tokenString := in.GetToken()

	_, err := models.VerifyToken(tokenString)
	if err != nil {
		return &pb.TokenResponse{Success: false}, err
	}

	return &pb.TokenResponse{Success: true}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
