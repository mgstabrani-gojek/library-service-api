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
	book, _ := bookController.Repository.FindBookByID(id)
	json.NewEncoder(w).Encode(book)
}
