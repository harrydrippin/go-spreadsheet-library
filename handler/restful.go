package handler

import (
	"encoding/json"
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
	e.POST("/api/borrow", h.Borrow)
	e.POST("/api/return", h.Return)
	e.POST("/api/extend", h.Extend)
	e.GET("/api/status", h.Status)
}

func (h *RESTfulHandler) Healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "OK", "timestamp": time.Now().String()})
}

func (h *RESTfulHandler) Search(c echo.Context) error {
	title := c.QueryParam("title")
	books, err := h.service.Search(title)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}

func (h *RESTfulHandler) Borrow(c echo.Context) error {
	params := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	title, borrower := params["title"], params["borrower"]

	books, err := h.service.Search(title)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(books) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}

	if len(books) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too many books found")
	}

	book := books[0]
	book, err = h.service.Borrow(book, borrower)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (h *RESTfulHandler) Return(c echo.Context) error {
	params := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	title, borrower := params["title"], params["borrower"]

	books, err := h.service.Search(title)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(books) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}

	if len(books) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too many books found")
	}

	book := books[0]
	if book.Borrower != borrower {
		return echo.NewHTTPError(http.StatusBadRequest, "Book not borrowed by that borrower")
	}

	book, err = h.service.Return(book, borrower)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (h *RESTfulHandler) Extend(c echo.Context) error {
	params := make(map[string]string)
	err := json.NewDecoder(c.Request().Body).Decode(&params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	title, borrower := params["title"], params["borrower"]
	books, err := h.service.Search(title)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(books) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}

	if len(books) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Too many books found")
	}

	book := books[0]
	if book.Borrower != borrower {
		return echo.NewHTTPError(http.StatusBadRequest, "Book not borrowed by that borrower")
	}

	book, err = h.service.Extend(book, borrower)
	if err != nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
	}

	return c.JSON(http.StatusOK, book)
}

func (h *RESTfulHandler) Status(c echo.Context) error {
	borrower := c.QueryParam("borrower")

	books, err := h.service.Status(borrower)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}
