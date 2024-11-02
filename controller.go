package main

import (
	"log"
	"net/http"
)

var config *Config

func main() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	http.HandleFunc("/query", handleQuery)
	log.Printf("Starting server on %v", config.ControllerPort)
	log.Fatal(http.ListenAndServe(config.ControllerPort, nil))
}
