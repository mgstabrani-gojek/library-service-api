package main

import (
	"database/sql"
	"gojek/library-service-api/internal/controller"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const port = "8080"

func main() {
	dbConnection := "host=localhost user=postgres dbname=library sslmode=disable"
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/ping", controller.HandlePingRequest)
	http.HandleFunc("/healthz", controller.HandleHealthCheckRequest)

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
