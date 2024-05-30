package controller

import (
	"encoding/json"
	"net/http"
)

func HandlePingRequest(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
	}{Message: "pong"}

	jsonInBytes, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
}
