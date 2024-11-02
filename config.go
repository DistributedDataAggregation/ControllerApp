package main

import (
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
		log.Println("Error loading .env file")
		return nil, err
	}

	executorAddresses := strings.Split(os.Getenv("EXECUTOR_ADDRESSES"), ",")
	controllerPort := os.Getenv("CONTROLLER_PORT")
	dataPath := os.Getenv("DATA_PATH")

	return &Config{
		ExecutorAddresses: executorAddresses,
		ControllerPort:    controllerPort,
		DataPath:          dataPath,
	}, nil
}
