package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ExecutorAddresses []string
	ControllerPort    string
	DataPath          string
	MainExecutorIdx   int
	SwaggerHost       string
	AllowedOrigin     string
	ExecutorsPort     int32
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	executorAddresses := strings.Split(os.Getenv("EXECUTOR_ADDRESSES"), ",")
	controllerPort := os.Getenv("CONTROLLER_PORT")
	dataPath := os.Getenv("DATA_PATH")
	mainExecutorIdx := os.Getenv("MAIN_EXECUTOR_IDX")
	swaggerHost := os.Getenv("SWAGGER_HOST")
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	executorsPort := os.Getenv("EXECUTOR_EXECUTOR_PORT")

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
	if mainExecutorIdx == "" {
		log.Println("Error: MAIN_EXECUTOR_IDX is missing")
		return nil, fmt.Errorf("missing required environment variable: MAIN_EXECUTOR_IDX")
	}
	if swaggerHost == "" {
		log.Println("Error: SWAGGER_HOST is missing")
		return nil, fmt.Errorf("missing required environment variable: SWAGGER_HOST")
	}
	if executorsPort == "" {
		log.Println("Error: EXECUTOR_EXECUTOR_PORT is missing")
		return nil, fmt.Errorf("missing required environment variable: EXECUTOR_EXECUTOR_PORT")
	}

	idx, err := strconv.Atoi(mainExecutorIdx)
	if err != nil {
		log.Printf("Error parsing main executor index string to int: %v", err)
		return &Config{}, err
	}
	if idx >= len(executorAddresses) || idx < 0 {
		log.Printf("Invalid index for main executor passed")
		return &Config{}, fmt.Errorf("invalid index for main executor passed")
	}

	port, err := strconv.ParseInt(executorsPort, 10, 32)
	if err != nil {
		log.Printf("Error parsing main executor port string to int: %v", err)
		return &Config{}, err
	}

	return &Config{
		ExecutorAddresses: executorAddresses,
		ControllerPort:    controllerPort,
		DataPath:          dataPath,
		MainExecutorIdx:   idx,
		SwaggerHost:       swaggerHost,
		AllowedOrigin:     allowedOrigin,
		ExecutorsPort:     int32(port),
	}, nil
}
