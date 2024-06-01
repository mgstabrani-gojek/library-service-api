package main

import (
	"database/sql"
	"gojek/library-service-api/internal/controller"
	"gojek/library-service-api/internal/repository"
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

	bookRepository := &repository.BookRepository{DB: db}
	bookController := &controller.BookController{Repository: bookRepository}

	http.HandleFunc("/ping", controller.HandlePingRequest)
	http.HandleFunc("/healthz", controller.HandleHealthCheckRequest)
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookController.GetAllBooks(w, r)
		case http.MethodPost:
			bookController.AddBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookController.GetBookByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
