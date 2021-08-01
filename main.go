package main

import (
	"github.com/labstack/echo/v4"

	handlers "github.com/harrydrippin/go-spreadsheet-library/handler"
	repositories "github.com/harrydrippin/go-spreadsheet-library/repository"
	services "github.com/harrydrippin/go-spreadsheet-library/service"
	utils "github.com/harrydrippin/go-spreadsheet-library/utils"
)

func main() {
	e := echo.New()

	repository := repositories.NewSpreadsheetRepository(*utils.NewConfig())
	service := services.NewLibraryService(repository)
	handler := handlers.NewRESTfulHandler(service)
	handler.RegisterRoutes(e)

	e.Logger.Debug("Starting server on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
