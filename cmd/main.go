package main

import (
	"gojek/library-service-api/internal/controller"
	"log"
	"net/http"
)

const port = "8080"

func main() {
	http.HandleFunc("/ping", controller.HandlePingRequest)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
