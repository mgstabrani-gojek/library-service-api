package domain

type Book struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Price         float64 `json:"price"`
	PublishedDate string  `json:"publishedDate"`
}
