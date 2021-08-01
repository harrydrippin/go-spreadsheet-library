package handler

import (
	"net/http"
	"time"

	services "github.com/harrydrippin/go-spreadsheet-library/service"
	"github.com/labstack/echo/v4"
)

type RESTfulHandler struct {
	service services.LibraryUsecase
}

func NewRESTfulHandler(service services.LibraryUsecase) *RESTfulHandler {
	return &RESTfulHandler{service: service}
}

func (h *RESTfulHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.Healthcheck)
	e.GET("/api/search", h.Search)
}

func (h *RESTfulHandler) Healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "OK", "timestamp": time.Now().String()})
}

func (h *RESTfulHandler) Search(c echo.Context) error {
	title := c.QueryParam("title")
	book, err := h.service.Search(title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, book)
}
