package service

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-demo/pbdata/greeterpb"
)

func TestServerSayHello(t *testing.T) {
	server := &Server{Name: "Greeter"}

	reply, err := server.SayHello(context.Background(), &greeterpb.HelloRequest{Name: "gRPC"})
	if err != nil {
		t.Fatalf("SayHello() error = %v", err)
	}
	if got, want := reply.GetMessage(), "Hello, gRPC! (Greeter)"; got != want {
		t.Fatalf("SayHello() = %q, want %q", got, want)
	}
}

func TestServerSayHelloRequiresName(t *testing.T) {
	server := &Server{Name: "Greeter"}

	_, err := server.SayHello(context.Background(), &greeterpb.HelloRequest{})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("SayHello() code = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
}
