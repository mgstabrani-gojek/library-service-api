package controller_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"gojek/library-service-api/internal/controller"
	"gojek/library-service-api/internal/domain"
	"gojek/library-service-api/internal/repository"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func setupTestController(t *testing.T) (*controller.BookController, func()) {
	db := setupTestDB(t)
	bookRepository := &repository.BookRepository{DB: db}
	controller := &controller.BookController{Repository: bookRepository}
	return controller, func() {
		db.Close()
	}
}

func TestGetAllBooks_GivenNothing_ThenReturnEmptyBooksResponse(t *testing.T) {
	bookController, teardown := setupTestController(t)
	defer teardown()

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()
	bookController.GetAllBooks(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	response := map[string][]domain.Book{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response["books"])
}

func TestGetBookById_GivenExistedBook_ThenReturnCorrespondingBookResponse(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	bookController, teardown := setupTestController(t)
	defer teardown()

	book := &domain.Book{Title: "Clean Code", Price: 10.99, PublishedDate: "1990-06-01"}
	db.QueryRow(
		"INSERT INTO books (title, price, published_date) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.Price, book.PublishedDate).Scan(&book.ID)
	req := httptest.NewRequest(http.MethodGet, "/books/"+strconv.Itoa(book.ID), nil)
	w := httptest.NewRecorder()
	bookController.GetBookByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	response := domain.Book{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Clean Code", response.Title)

	db.Exec("DELETE FROM books WHERE id = $1", book.ID)
}

func TestGetBookById_GivenNotFoundBook_ThenReturnErrorResponse(t *testing.T) {
	bookController, teardown := setupTestController(t)
	defer teardown()

	req := httptest.NewRequest(http.MethodGet, "/books/"+strconv.Itoa(-1), nil)
	w := httptest.NewRecorder()
	bookController.GetBookByID(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	expectedResponse := `{"error":"Internal server error."}`

	assert.Equal(t, expectedResponse, string(data))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestAddBook_GivenInvalidRequestBody_ThenReturnErrorResponse(t *testing.T) {
	bookController, teardown := setupTestController(t)
	defer teardown()

	invalidRequest := struct {
		ID string `json:"id"`
	}{ID: "invalid request"}
	invalidRequestInJSON, _ := json.Marshal(invalidRequest)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(invalidRequestInJSON))
	w := httptest.NewRecorder()
	bookController.AddBook(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	expectedResponse := `{"error":"Internal server error."}`

	assert.Equal(t, expectedResponse, string(data))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestAddBook_GivenValidRequestBody_ThenReturnAddedBookResponse(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	bookController, teardown := setupTestController(t)
	defer teardown()

	book := &domain.Book{Title: "The Great Gatsby", Price: 15.99, PublishedDate: "1925-04-10"}
	bookJSON, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(bookJSON))
	w := httptest.NewRecorder()
	bookController.AddBook(w, req)

	res := w.Result()
	defer res.Body.Close()

	response := map[string]interface{}{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "The Great Gatsby", response["title"])
	assert.Equal(t, "Book successfully added to the library.", response["message"])
	assert.Equal(t, http.StatusOK, res.StatusCode)

	db.Exec("DELETE FROM books WHERE id = $1", response["id"])
}

func TestUpdateBookTitle_GivenInvalidRequestBody_ThenReturnErrorResponse(t *testing.T) {
	bookController, teardown := setupTestController(t)
	defer teardown()

	invalidRequest := struct {
		Title int `json:"title"`
	}{Title: 1234}
	invalidRequestInJSON, _ := json.Marshal(invalidRequest)
	req := httptest.NewRequest(http.MethodPut, "/books", bytes.NewReader(invalidRequestInJSON))
	w := httptest.NewRecorder()
	bookController.UpdateBookTitle(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	expectedResponse := `{"error":"Internal server error."}`

	assert.Equal(t, expectedResponse, string(data))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUpdateBoookTitle_GivenNotFoundBook_ThenReturnErrorResponse(t *testing.T) {
	bookController, teardown := setupTestController(t)
	defer teardown()

	bodyRequest := struct {
		Title string `json:"title"`
	}{Title: "Updated Title"}
	bodyRequestInJSON, _ := json.Marshal(bodyRequest)
	req := httptest.NewRequest(http.MethodGet, "/books/"+strconv.Itoa(-1), bytes.NewReader(bodyRequestInJSON))
	w := httptest.NewRecorder()
	bookController.UpdateBookTitle(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	expectedResponse := `{"error":"Internal server error."}`

	assert.Equal(t, expectedResponse, string(data))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUpdateBookTitle_GivenValidRequestBodyAndExistedBook_ThenReturnUpdatedBookResponse(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	bookController, teardown := setupTestController(t)
	defer teardown()

	book := &domain.Book{Title: "Clean Code", Price: 10.99, PublishedDate: "1990-06-01"}
	db.QueryRow(
		"INSERT INTO books (title, price, published_date) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.Price, book.PublishedDate).Scan(&book.ID)

	bookTitleUpdate := struct {
		Title string `json:"title"`
	}{Title: "Updated Title"}
	bookJSON, _ := json.Marshal(bookTitleUpdate)
	req := httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(book.ID), bytes.NewReader(bookJSON))
	w := httptest.NewRecorder()
	bookController.UpdateBookTitle(w, req)

	res := w.Result()
	defer res.Body.Close()

	response := map[string]interface{}{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", response["title"])
	assert.Equal(t, "Book title successfully updated.", response["message"])
	assert.Equal(t, http.StatusOK, res.StatusCode)

	db.Exec("DELETE FROM books WHERE id = $1", response["id"])
}

func TestDeleteBookById_GivenNotFoundBook_ThenReturnErrorResponse(t *testing.T) {
	bookController, teardown := setupTestController(t)
	defer teardown()

	req := httptest.NewRequest(http.MethodDelete, "/books/"+strconv.Itoa(-1), nil)
	w := httptest.NewRecorder()
	bookController.DeleteBookByID(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	expectedResponse := `{"error":"Internal server error."}`

	assert.Equal(t, expectedResponse, string(data))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestDeleteBookById_GivenExistedBook_ThenReturnDeletedBookResponse(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	bookController, teardown := setupTestController(t)
	defer teardown()

	book := &domain.Book{Title: "Clean Code", Price: 10.99, PublishedDate: "1990-06-01"}
	db.QueryRow(
		"INSERT INTO books (title, price, published_date) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.Price, book.PublishedDate).Scan(&book.ID)

	req := httptest.NewRequest(http.MethodDelete, "/books/"+strconv.Itoa(book.ID), nil)
	w := httptest.NewRecorder()
	bookController.DeleteBookByID(w, req)

	res := w.Result()
	defer res.Body.Close()

	response := map[string]interface{}{}
	err := json.NewDecoder(res.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Book successfully deleted.", response["message"])
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
