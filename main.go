package main

import (
	"github.com/labstack/echo/v4"

	model "github.com/harrydrippin/scatterlab-library/model"
	repositories "github.com/harrydrippin/scatterlab-library/repository"
	utils "github.com/harrydrippin/scatterlab-library/utils"
)

func main() {
	e := echo.New()

	// Set up the repositories
	repository := repositories.NewSpreadsheetRepository(*utils.NewConfig())

	books, err := repository.GetByTitleSubstring("에픽테토스")
	if err != nil {
		panic(err)
	}

	book := books[0]
	book.Status = model.StatusBorrowed
	err = repository.Update(book)
	if err != nil {
		panic(err)
	}

	e.Logger.Debug("Starting server on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
