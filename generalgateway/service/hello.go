package service

import (
	"context"
	"time"

	"grpc-demo/generalgateway/grpc_client"
	"grpc-demo/pbdata/greeterpb"
)

func SayHello(ctx context.Context, name string) (*greeterpb.HelloReply, error) {
	callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return grpc_client.GreeterClient.SayHello(callCtx, &greeterpb.HelloRequest{Name: name})
}
