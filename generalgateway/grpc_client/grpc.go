package grpc_client

import (
	"errors"
	"fmt"

	"grpc-demo/pbdata/calculatorpb"
	"grpc-demo/pbdata/greeterpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GreeterClient    greeterpb.GreeterClient
	CalculatorClient calculatorpb.CalculatorClient
	greeterConn      *grpc.ClientConn
	calculatorConn   *grpc.ClientConn
)

type RpcClientConfig struct {
	GreeterAddress    string `yaml:"greeter_address"`
	CalculatorAddress string `yaml:"calculator_address"`
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

func NewCalculatorClient(address string) (calculatorpb.CalculatorClient, *grpc.ClientConn, error) {
	fmt.Println("calculator address:", address)
	conn, err := newCommonConn(address)
	if err != nil {
		return nil, nil, err
	}
	return calculatorpb.NewCalculatorClient(conn), conn, nil
}

func Init(conf RpcClientConfig) error {
	var err error
	GreeterClient, greeterConn, err = NewGreeterClient(conf.GreeterAddress)
	if err != nil {
		return err
	}

	CalculatorClient, calculatorConn, err = NewCalculatorClient(conf.CalculatorAddress)
	if err != nil {
		_ = greeterConn.Close()
		GreeterClient = nil
		greeterConn = nil
		return err
	}
	return nil
}

func Close() error {
	var errs []error
	if greeterConn != nil {
		errs = append(errs, greeterConn.Close())
		greeterConn = nil
	}
	if calculatorConn != nil {
		errs = append(errs, calculatorConn.Close())
		calculatorConn = nil
	}
	return errors.Join(errs...)
}
