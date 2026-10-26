package service

import (
	"context"
	"time"

	"grpc-demo/generalgateway/grpc_client"
	"grpc-demo/pbdata/calculatorpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	OperationAdd      = "add"
	OperationMultiply = "multiply"
)

func Calculate(ctx context.Context, operation string, left, right int64) (*calculatorpb.OperationReply, error) {
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return calculate(callCtx, grpc_client.CalculatorClient, operation, left, right)
}

func calculate(ctx context.Context, client calculatorpb.CalculatorClient, operation string, left, right int64) (*calculatorpb.OperationReply, error) {
	req := &calculatorpb.BinaryOperationRequest{Left: left, Right: right}
	switch operation {
	case OperationAdd:
		return client.Add(ctx, req)
	case OperationMultiply:
		return client.Multiply(ctx, req)
	default:
		return nil, status.Error(codes.InvalidArgument, "unsupported operation")
	}
}
