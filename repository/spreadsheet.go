package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	models "github.com/harrydrippin/go-spreadsheet-library/model"
	utils "github.com/harrydrippin/go-spreadsheet-library/utils"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// BookRepository is a repository for a book
type BookRepository interface {
	GetByTitleSubstring(title string) ([]models.Book, error)
	GetAll() ([]models.Book, error)
	Update(book models.Book) error
}

type SpreadsheetRepository struct {
	config       utils.Config
	client       *http.Client
	sheetService *sheets.Service
}

func NewSpreadsheetRepository(config utils.Config) *SpreadsheetRepository {
	ctx := context.Background()
	client := utils.GetGoogleClient(config.GoogleCredentialJSON)
	sheetService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Unable to retrieve Sheets client", err)
	}

	return &SpreadsheetRepository{
		config:       config,
		client:       client,
		sheetService: sheetService,
	}
}

func (s *SpreadsheetRepository) GetByTitleSubstring(title string) ([]models.Book, error) {
	books, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	result := []models.Book{}
	for _, book := range books {
		if strings.Contains(book.Title, title) {
			result = append(result, book)
		}
	}

	return result, nil
}

func (r *SpreadsheetRepository) GetAll() ([]models.Book, error) {
	readRange := fmt.Sprintf("%s!A3:H", r.config.GoogleSpreadsheetName)
	response, err := r.sheetService.Spreadsheets.Values.Get(r.config.GoogleSpreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	books := make([]models.Book, 0)
	for _, row := range response.Values {
		book := models.Book{}
		bookId, err := strconv.ParseInt(row[0].(string), 0, 64)
		if err != nil {
			return nil, err
		}
		book.ID = int(bookId)
		book.Title = row[1].(string)
		book.Author = row[2].(string)
		book.Publisher = row[3].(string)
		book.Position = row[4].(string)
		book.Status = row[5].(string)
		if len(row) == 6 {
			book.Borrower = ""
			book.DueDate = ""
		} else {
			book.Borrower = row[6].(string)
			book.DueDate = row[7].(string)
		}

		books = append(books, book)
	}

	return books, nil
}

func (r *SpreadsheetRepository) Update(book models.Book) error {
	rowId := book.ID + 2
	readRange := fmt.Sprintf("%s!A%d:H%d", r.config.GoogleSpreadsheetName, rowId, rowId)
	_, err := r.sheetService.Spreadsheets.Values.Get(r.config.GoogleSpreadsheetID, readRange).Do()
	if err != nil {
		return err
	}

	valueRange := sheets.ValueRange{
		Values: [][]interface{}{
			{book.ID, book.Title, book.Author, book.Publisher, book.Position, book.Status, book.Borrower, book.DueDate},
		},
	}

	call := r.sheetService.Spreadsheets.Values.Update(r.config.GoogleSpreadsheetID, readRange, &valueRange).ValueInputOption("RAW")
	_, err = call.Do()
	if err != nil {
		return err
	}

	return nil
}
