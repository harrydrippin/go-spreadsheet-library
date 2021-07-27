package main

import (
	"fmt"
	"github.com/labstack/echo/v4"

	repositories "github.com/harrydrippin/scatterlab-library/repository"
	utils "github.com/harrydrippin/scatterlab-library/utils"
)

func main() {
	e := echo.New()

	// Set up the repositories
	repository := repositories.NewSpreadsheetRepository(*utils.NewConfig())

	books, err := repository.GetAll()
	if err != nil {
		panic(err)
	}

	// print books
	for _, book := range books {
		fmt.Printf("%s: %s (%s)\n", book.Title, book.Author, book.Publisher)
	}

	e.Logger.Debug("Starting server on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
