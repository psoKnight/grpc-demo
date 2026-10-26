package grpc_client

import (
	"fmt"

	"grpc-demo/pbdata/greeterpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GreeterClient greeterpb.GreeterClient
	greeterConn   *grpc.ClientConn
)

type RpcClientConfig struct {
	GreeterAddress string `yaml:"greeter_address"`
}

func newCommonConn(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*20)),
	)
}

func NewGreeterClient(address string) (greeterpb.GreeterClient, *grpc.ClientConn, error) {
	fmt.Println("greeter address:", address)
	conn, err := newCommonConn(address)
	if err != nil {
		return nil, nil, err
	}
	return greeterpb.NewGreeterClient(conn), conn, nil
}

func Init(conf RpcClientConfig) error {
	var err error
	GreeterClient, greeterConn, err = NewGreeterClient(conf.GreeterAddress)
	return err
}

func Close() error {
	if greeterConn == nil {
		return nil
	}
	return greeterConn.Close()
}
