package service

import (
	"errors"
	"time"

	model "github.com/harrydrippin/go-spreadsheet-library/model"
	repositories "github.com/harrydrippin/go-spreadsheet-library/repository"
)

// LibraryUsecase is the interface that defines the usecase for the library
type LibraryUsecase interface {
	Search(title string) ([]model.Book, error)
	Borrow(book model.Book, borrower string) error
	Return(book model.Book) error
	Extend(book model.Book) error
}

// LibraryService is the service that handles the library usecase
type LibraryService struct {
	repository repositories.BookRepository
}

// NewLibraryService returns a new instance of LibraryService
func NewLibraryService(repository repositories.BookRepository) *LibraryService {
	return &LibraryService{repository: repository}
}

func (library *LibraryService) Search(title string) ([]model.Book, error) {
	return library.repository.GetByTitleSubstring(title)
}

func (library *LibraryService) Borrow(book model.Book, borrower string) error {
	if book.Status != model.StatusInOffice {
		return errors.New("This book is already borrowed")
	}

	// 4주 뒤 반납 예정인 대출된 책으로 변경
	book.Status = model.StatusBorrowed
	book.Borrower = borrower
	book.DueDate = time.Now().AddDate(0, 0, 28).Format("2006-01-02")
	return library.repository.Update(book)
}

func (library *LibraryService) Return(book model.Book) error {
	if book.Status == model.StatusInOffice {
		return errors.New("This book is not borrowed")
	}

	// 반납된 책으로 변경
	book.Status = model.StatusInOffice
	book.Borrower = ""
	book.DueDate = ""
	return library.repository.Update(book)
}

func (library *LibraryService) Extend(book model.Book) error {
	if book.Status != model.StatusBorrowed {
		return errors.New("This book is not borrowed or delayed")
	}

	// 대출 기한을 연장
	book.DueDate = time.Now().AddDate(0, 0, 28).Format("2006-01-02")
	return library.repository.Update(book)
}
