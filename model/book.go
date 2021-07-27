package model

// Constants for representing book status
const (
	StatusInOffice = "사내 비치"
	StatusBorrowed = "대출"
	StatusOverdue  = "연체"
)

// Book is a simple data structure to represent a book.
type Book struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Position  string `json:"position"`
	Status    string `json:"status"`
	Borrower  string `json:"borrower"`
	DueDate   string `json:"due_date"`
}

// NewBook creates a new book with the given parameters.
func NewBook(id int, title, author, publisher, position, status, borrower, dueDate string) Book {
	return Book{
		ID:        id,
		Title:     title,
		Author:    author,
		Publisher: publisher,
		Position:  position,
		Status:    status,
		Borrower:  borrower,
		DueDate:   dueDate,
	}
}
