package config

import (
	"fmt"
	"os"

	"grpc-demo/generalgateway/grpc_client"

	"gopkg.in/yaml.v3"
)

var Conf *Config

type ServerConfig struct {
	Address string `yaml:"address"`
}

type GinConfig struct {
	Mode string `yaml:"mode"`
}

type Config struct {
	Server           ServerConfig                `yaml:"server"`
	Gin              GinConfig                   `yaml:"gin"`
	GRpcClientConfig grpc_client.RpcClientConfig `yaml:"rpcclient"`
}

func Init(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}
	if conf.Server.Address == "" {
		return fmt.Errorf("server.address is required")
	}
	if conf.GRpcClientConfig.GreeterAddress == "" {
		return fmt.Errorf("rpcclient.greeter_address is required")
	}
	if conf.GRpcClientConfig.CalculatorAddress == "" {
		return fmt.Errorf("rpcclient.calculator_address is required")
	}

	Conf = &conf
	return nil
}
