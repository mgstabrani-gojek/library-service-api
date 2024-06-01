package controller

import (
	"encoding/json"
	"net/http"
)

func HandleHealthCheckRequest(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
	}{Message: "service is available"}

	jsonInBytes, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
}
