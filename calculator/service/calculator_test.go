package service

import (
	"context"
	"math"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-demo/pbdata/calculatorpb"
)

func TestServerAdd(t *testing.T) {
	server := &Server{Name: "Calculator"}

	reply, err := server.Add(context.Background(), &calculatorpb.BinaryOperationRequest{
		Left:  21,
		Right: 21,
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if got, want := reply.GetResult(), int64(42); got != want {
		t.Fatalf("Add() = %d, want %d", got, want)
	}
}

func TestServerMultiply(t *testing.T) {
	server := &Server{Name: "Calculator"}

	reply, err := server.Multiply(context.Background(), &calculatorpb.BinaryOperationRequest{
		Left:  6,
		Right: 7,
	})
	if err != nil {
		t.Fatalf("Multiply() error = %v", err)
	}
	if got, want := reply.GetResult(), int64(42); got != want {
		t.Fatalf("Multiply() = %d, want %d", got, want)
	}
}

func TestServerRejectsOverflow(t *testing.T) {
	server := &Server{Name: "Calculator"}

	t.Run("add", func(t *testing.T) {
		_, err := server.Add(context.Background(), &calculatorpb.BinaryOperationRequest{
			Left:  math.MaxInt64,
			Right: 1,
		})
		if status.Code(err) != codes.OutOfRange {
			t.Fatalf("Add() code = %v, want %v", status.Code(err), codes.OutOfRange)
		}
	})

	t.Run("multiply", func(t *testing.T) {
		_, err := server.Multiply(context.Background(), &calculatorpb.BinaryOperationRequest{
			Left:  math.MaxInt64,
			Right: 2,
		})
		if status.Code(err) != codes.OutOfRange {
			t.Fatalf("Multiply() code = %v, want %v", status.Code(err), codes.OutOfRange)
		}
	})
}
