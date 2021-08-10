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
	config := utils.NewConfig()

	repository := repositories.NewSpreadsheetRepository(*config)
	service := services.NewLibraryService(repository)

	restfulHandler := handlers.NewRESTfulHandler(service)
	restfulHandler.RegisterRoutes(e)
	slackHandler := handlers.NewSlackHandler(service, *config)
	slackHandler.RegisterRoutes(e)

	e.Logger.Debug("Starting server on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
