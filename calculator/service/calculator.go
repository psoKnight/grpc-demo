package service

import (
	"context"
	"math"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-demo/pbdata/calculatorpb"
)

type Server struct {
	calculatorpb.UnimplementedCalculatorServer
	Name string
}

func (s *Server) Add(_ context.Context, req *calculatorpb.BinaryOperationRequest) (*calculatorpb.OperationReply, error) {
	left, right := req.GetLeft(), req.GetRight()
	if (right > 0 && left > math.MaxInt64-right) ||
		(right < 0 && left < math.MinInt64-right) {
		return nil, status.Error(codes.OutOfRange, "integer overflow")
	}

	return &calculatorpb.OperationReply{Result: left + right}, nil
}

func (s *Server) Multiply(_ context.Context, req *calculatorpb.BinaryOperationRequest) (*calculatorpb.OperationReply, error) {
	result := new(big.Int).Mul(big.NewInt(req.GetLeft()), big.NewInt(req.GetRight()))
	if !result.IsInt64() {
		return nil, status.Error(codes.OutOfRange, "integer overflow")
	}

	return &calculatorpb.OperationReply{Result: result.Int64()}, nil
}
