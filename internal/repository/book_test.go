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

func TestGetAllBooks_GivenNothing_ThenReturnEmptyBooks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	bookRepository := &repository.BookRepository{DB: db}
	books, err := bookRepository.FindAllBooks()
	assert.NoError(t, err)
	assert.IsType(t, []domain.Book{}, books)
}
