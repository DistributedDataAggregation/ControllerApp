package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ExecutorAddresses []string
	ControllerPort    string
	DataPath          string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	executorAddresses := strings.Split(os.Getenv("EXECUTOR_ADDRESSES"), ",")
	controllerPort := os.Getenv("CONTROLLER_PORT")
	dataPath := os.Getenv("DATA_PATH")

	if len(executorAddresses) == 0 || executorAddresses[0] == "" {
		log.Println("Error: EXECUTOR_ADDRESSES is missing")
		return nil, fmt.Errorf("missing required environment variable: EXECUTOR_ADDRESSES")
	}
	if controllerPort == "" {
		log.Println("Error: CONTROLLER_PORT is missing")
		return nil, fmt.Errorf("missing required environment variable: CONTROLLER_PORT")
	}
	if dataPath == "" {
		log.Println("Error: DATA_PATH is missing")
		return nil, fmt.Errorf("missing required environment variable: DATA_PATH")
	}

	return &Config{
		ExecutorAddresses: executorAddresses,
		ControllerPort:    controllerPort,
		DataPath:          dataPath,
	}, nil
}
