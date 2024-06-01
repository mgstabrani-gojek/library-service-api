package controller

import (
	"encoding/json"
	"gojek/library-service-api/internal/domain"
	"gojek/library-service-api/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type BookController struct {
	Repository *repository.BookRepository
}

func (bookController *BookController) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, _ := bookController.Repository.FindAllBooks()
	json.NewEncoder(w).Encode(map[string][]domain.Book{"books": books})
}

func (bookController *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/books/"))
	book, err := bookController.Repository.FindBookByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := struct {
			Error string `json:"error"`
		}{Error: "Internal server error."}

		jsonInBytes, _ := json.Marshal(response)
		w.Write(jsonInBytes)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (bookController *BookController) AddBook(w http.ResponseWriter, r *http.Request) {
	book := domain.Book{}
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := struct {
			Error string `json:"error"`
		}{Error: "Internal server error."}

		jsonInBytes, _ := json.Marshal(response)
		w.Write(jsonInBytes)
		return
	}

	bookController.Repository.SaveBook(&book)
	bookResponse := map[string]interface{}{
		"id":            book.ID,
		"title":         book.Title,
		"price":         book.Price,
		"publishedDate": book.PublishedDate,
		"message":       "Book successfully added to the library.",
	}
	json.NewEncoder(w).Encode(bookResponse)
}
