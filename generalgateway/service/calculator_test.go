package service

import (
	"context"
	"testing"

	"grpc-demo/pbdata/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeCalculatorClient struct {
	addCalled      bool
	multiplyCalled bool
	request        *calculatorpb.BinaryOperationRequest
}

func (f *fakeCalculatorClient) Add(_ context.Context, req *calculatorpb.BinaryOperationRequest, _ ...grpc.CallOption) (*calculatorpb.OperationReply, error) {
	f.addCalled = true
	f.request = req
	return &calculatorpb.OperationReply{Result: req.GetLeft() + req.GetRight()}, nil
}

func (f *fakeCalculatorClient) Multiply(_ context.Context, req *calculatorpb.BinaryOperationRequest, _ ...grpc.CallOption) (*calculatorpb.OperationReply, error) {
	f.multiplyCalled = true
	f.request = req
	return &calculatorpb.OperationReply{Result: req.GetLeft() * req.GetRight()}, nil
}

func TestCalculateDispatchesOperation(t *testing.T) {
	tests := []struct {
		name         string
		operation    string
		left         int64
		right        int64
		want         int64
		wantAdd      bool
		wantMultiply bool
	}{
		{name: "add", operation: OperationAdd, left: 21, right: 21, want: 42, wantAdd: true},
		{name: "multiply", operation: OperationMultiply, left: 21, right: 2, want: 42, wantMultiply: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &fakeCalculatorClient{}
			reply, err := calculate(context.Background(), client, tt.operation, tt.left, tt.right)
			if err != nil {
				t.Fatalf("calculate() error = %v", err)
			}
			if reply.GetResult() != tt.want {
				t.Fatalf("calculate() = %d, want %d", reply.GetResult(), tt.want)
			}
			if client.addCalled != tt.wantAdd || client.multiplyCalled != tt.wantMultiply {
				t.Fatalf("dispatched add=%v multiply=%v, want add=%v multiply=%v", client.addCalled, client.multiplyCalled, tt.wantAdd, tt.wantMultiply)
			}
			if got := client.request; got.GetLeft() != tt.left || got.GetRight() != tt.right {
				t.Fatalf("request = %+v, want left=%d right=%d", got, tt.left, tt.right)
			}
		})
	}
}

func TestCalculateRejectsUnknownOperation(t *testing.T) {
	client := &fakeCalculatorClient{}
	_, err := calculate(context.Background(), client, "divide", 21, 2)
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("calculate() code = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
	if client.addCalled || client.multiplyCalled {
		t.Fatal("unknown operation must not call downstream")
	}
}
