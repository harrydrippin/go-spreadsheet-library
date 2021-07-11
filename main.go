package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Logger.Debug("Starting server on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
