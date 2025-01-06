package main

import (
	"encoding/json"
	"net/http"
)

// @Summary Health check endpoint
// @Description Checks controller status
// @Tags health check
// @Produce  json
// @Success 200 {string} string "Health check passed"
// @Router /status [get]
func handleStatusCheck(w http.ResponseWriter, r *http.Request) {
	result := "Controller is ready"
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
