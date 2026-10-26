package service

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-demo/pbdata/greeterpb"
)

type Server struct {
	greeterpb.UnimplementedGreeterServer
	Name string
}

func (s *Server) SayHello(_ context.Context, req *greeterpb.HelloRequest) (*greeterpb.HelloReply, error) {
	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	return &greeterpb.HelloReply{
		Message: fmt.Sprintf("Hello, %s! (%s)", name, s.Name),
	}, nil
}
