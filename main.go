package main

import (
	"log"
	"net/http"

	"controller/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

var config *Config

// @title           Swagger Distributed data aggregation system API
// @version         1.0
// @BasePath  /api/v1
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	executorsClient := NewExecutorsClient()
	executorsClient.OpenSockets()
	planner := NewPlanner()
	processor := NewProcessor(planner, executorsClient)
	scheduler := NewQueriesScheduler(processor)
	queryHandler := NewQueryHandler(scheduler)

	docs.SwaggerInfo.Host = config.SwaggerHost
	http.HandleFunc("/api/v1/query", queryHandler.handleQuery)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Printf("Starting server on %v", config.ControllerPort)
	log.Fatal(http.ListenAndServe(config.ControllerPort, nil))
}
