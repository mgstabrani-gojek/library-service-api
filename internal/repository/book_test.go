package repository_test

import (
	"database/sql"
	"gojek/library-service-api/internal/domain"
	"gojek/library-service-api/internal/repository"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost user=postgres dbname=library sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS books (
            id SERIAL PRIMARY KEY,
            title VARCHAR(100),
            price NUMERIC(10, 2),
            published_date DATE
        );
    `)
	if err != nil {
		t.Fatalf("Failed to create books table: %v", err)
	}

	return db
}

func TestFindAllBooks_GivenNothing_ThenReturnEmptyBooks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	bookRepository := &repository.BookRepository{DB: db}
	books, err := bookRepository.FindAllBooks()
	assert.NoError(t, err)
	assert.IsType(t, []domain.Book{}, books)
}

func TestFindAllBooks_GivenOneBook_ThenReturnListOfBooks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	book := &domain.Book{ID: 1, Title: "Clean Code", Price: 15.99, PublishedDate: "1990-06-01T00:00:00Z"}
	db.QueryRow(
		"INSERT INTO books (title, price, published_date) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.Price, book.PublishedDate).Scan(&book.ID)

	bookRepository := &repository.BookRepository{DB: db}
	books, _ := bookRepository.FindAllBooks()
	assert.Equal(t, []domain.Book{*book}, books)

	db.Exec("DELETE FROM books WHERE id = $1", book.ID)
}

func TestFindBookById_GivenNotFoundBook_ThenReturnError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	bookRepository := &repository.BookRepository{DB: db}
	_, err := bookRepository.FindBookByID(1)
	assert.Error(t, err, sql.ErrNoRows)
}

func TestFindBookById_GivenExistedBook_ThenReturnCorrespondingBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	createdBook := domain.Book{Title: "Clean Code", Price: 15.99, PublishedDate: "1990-06-01T00:00:00Z"}
	db.QueryRow(
		"INSERT INTO books (title, price, published_date) VALUES ($1, $2, $3) RETURNING id",
		createdBook.Title, createdBook.Price, createdBook.PublishedDate).Scan(&createdBook.ID)

	bookRepository := &repository.BookRepository{DB: db}
	book, _ := bookRepository.FindBookByID(createdBook.ID)
	assert.Equal(t, createdBook, book)

	db.Exec("DELETE FROM books WHERE id = $1", book.ID)
}

func TestSaveBook_GivenNewBook_ThenNewBookInserted(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	book := &domain.Book{Title: "Clean Code", Price: 15.99, PublishedDate: "1990-06-01T00:00:00Z"}
	bookRepository := &repository.BookRepository{DB: db}
	err := bookRepository.SaveBook(book)
	assert.NoError(t, err)

	db.Exec("DELETE FROM books WHERE id = $1", book.ID)
}
