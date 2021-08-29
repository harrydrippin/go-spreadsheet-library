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
	Borrow(book model.Book, borrower string) (model.Book, error)
	Return(book model.Book, borrower string) (model.Book, error)
	Extend(book model.Book, borrower string) (model.Book, error)
	Status(borrower string) ([]model.Book, error)
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
	return library.repository.Search(title)
}

func (library *LibraryService) Borrow(book model.Book, borrower string) (model.Book, error) {
	if book.Status != model.StatusInOffice {
		return model.Book{}, errors.New("이미 대출되어 있는 책이에요. 다시 확인해주세요.")
	}

	// 4주 뒤 반납 예정인 대출된 책으로 변경
	book.Status = model.StatusBorrowed
	book.Borrower = borrower
	book.DueDate = time.Now().AddDate(0, 0, 28).Format("2006-01-02")
	err := library.repository.Update(book)

	return book, err
}

func (library *LibraryService) Return(book model.Book, borrower string) (model.Book, error) {
	if book.Status == model.StatusInOffice {
		return model.Book{}, errors.New("대출된 책이 아니에요! 다시 확인해주세요.")
	}

	if book.Borrower != borrower {
		return model.Book{}, errors.New("@" + borrower + " 님이 대출하신 책이 아니에요. 다시 확인해주세요!")
	}

	// 반납된 책으로 변경
	book.Status = model.StatusInOffice
	book.Borrower = ""
	book.DueDate = ""
	err := library.repository.Update(book)

	return book, err
}

func (library *LibraryService) Extend(book model.Book, borrower string) (model.Book, error) {
	if book.Status != model.StatusBorrowed {
		return model.Book{}, errors.New("이 책은 대출된 상태가 아니에요. 다시 확인해주세요.")
	}

	if book.Borrower != borrower {
		return model.Book{}, errors.New("@" + borrower + " 님이 대출하신 책이 아니에요. 다시 확인해주세요!")
	}

	// 대출 기한을 연장
	book.DueDate = time.Now().AddDate(0, 0, 28).Format("2006-01-02")
	err := library.repository.Update(book)

	return book, err
}

func (library *LibraryService) Status(borrower string) ([]model.Book, error) {
	books, err := library.repository.GetAll()
	if err != nil {
		return nil, err
	}

	result := []model.Book{}
	for _, book := range books {
		if book.Borrower == borrower {
			result = append(result, book)
		}
	}

	return result, nil
}
