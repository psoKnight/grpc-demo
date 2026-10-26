package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Conf *Config

type RPCConfig struct {
	Protocol string `yaml:"protocol"`
	Address  string `yaml:"address"`
}

type HealthConfig struct {
	Address string `yaml:"address"`
}

type Config struct {
	RPCServer    RPCConfig    `yaml:"rpcserver"`
	HealthServer HealthConfig `yaml:"healthserver"`
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
	if conf.RPCServer.Protocol == "" {
		return fmt.Errorf("rpcserver.protocol is required")
	}
	if conf.RPCServer.Address == "" {
		return fmt.Errorf("rpcserver.address is required")
	}
	if conf.HealthServer.Address == "" {
		return fmt.Errorf("healthserver.address is required")
	}

	Conf = &conf
	return nil
}
