package repository

import (
	"database/sql"
	"gojek/library-service-api/internal/domain"
)

type BookRepository struct {
	DB *sql.DB
}

func (bookRepository *BookRepository) FindAllBooks() ([]domain.Book, error) {
	rows, _ := bookRepository.DB.Query("SELECT id, title, price, published_date FROM books")
	defer rows.Close()

	books := []domain.Book{}
	return books, nil
}
