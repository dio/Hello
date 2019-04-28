package hello

import (
	// stdlib
	"context"

	// project
	pb "github.com/dio/hello/api/hello/v1"
)

type service struct{}

// NewService ...
func NewService() pb.HelloServiceServer {
	return &service{}
}

func (s *service) Say(ctx context.Context, req *pb.SayRequest) (*pb.SayResponse, error) {
	return &pb.SayResponse{
		Message: "Hello, " + req.Message,
	}, nil
}

func (s *service) Poke(ctx context.Context, req *pb.PokeRequest) (*pb.PokeResponse, error) {
	return &pb.PokeResponse{
		Message: "Hi, " + req.Name,
	}, nil
}
