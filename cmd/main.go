package main

import (
	"database/sql"
	"gojek/library-service-api/internal/config"
	"gojek/library-service-api/internal/controller"
	"gojek/library-service-api/internal/repository"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	port := config.GetEnv("PORT", "8080")

	dbConfig := config.NewDBConfig()

	dbConnection := "host=" + dbConfig.Host + " user=" + dbConfig.User + " dbname=" + dbConfig.DBName + " sslmode=" + dbConfig.SSLMode
	if dbConfig.Password != "" {
		dbConnection += " password=" + dbConfig.Password
	}

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
		case http.MethodPut:
			bookController.UpdateBookTitle(w, r)
		case http.MethodDelete:
			bookController.DeleteBookByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
