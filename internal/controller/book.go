package controller

import (
	"encoding/json"
	"gojek/library-service-api/internal/domain"
	"gojek/library-service-api/internal/repository"
	"net/http"
)

type BookController struct {
	Repository *repository.BookRepository
}

func (bookController *BookController) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, _ := bookController.Repository.FindAllBooks()
	json.NewEncoder(w).Encode(map[string][]domain.Book{"books": books})
}
