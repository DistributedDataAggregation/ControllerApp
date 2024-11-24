package main

import (
	"os"
	"strings"
)

type Config struct {
	ExecutorAddresses []string
	ControllerPort    string
	DataPath          string
}

func LoadConfig() (*Config, error) {
	executorAddresses := strings.Split(os.Getenv("EXECUTOR_ADDRESSES"), ",")
	controllerPort := os.Getenv("CONTROLLER_PORT")
	dataPath := os.Getenv("DATA_PATH")

	return &Config{
		ExecutorAddresses: executorAddresses,
		ControllerPort:    controllerPort,
		DataPath:          dataPath,
	}, nil
}
