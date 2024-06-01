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
	for rows.Next() {
		book := domain.Book{}
		rows.Scan(&book.ID, &book.Title, &book.Price, &book.PublishedDate)
		books = append(books, book)
	}
	return books, nil
}

func (bookRepository *BookRepository) FindBookByID(id int) (domain.Book, error) {
	book := domain.Book{}
	err := bookRepository.DB.QueryRow("SELECT id, title, price, published_date FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Price, &book.PublishedDate)
	return book, err
}

func (bookRepository *BookRepository) SaveBook(book *domain.Book) error {
	return bookRepository.DB.QueryRow(
		"INSERT INTO books (title, price, published_date) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.Price, book.PublishedDate).Scan(&book.ID)
}
