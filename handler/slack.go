package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	services "github.com/harrydrippin/go-spreadsheet-library/service"
	utils "github.com/harrydrippin/go-spreadsheet-library/utils"
	views "github.com/harrydrippin/go-spreadsheet-library/view"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
)

type SlackHandler struct {
	SigningSecret string

	client  *slack.Client
	userId  string
	service services.LibraryUsecase
}

func NewSlackHandler(service services.LibraryUsecase, config utils.Config) *SlackHandler {
	client := slack.New(config.SlackToken)
	bot, err := client.AuthTest()
	if err != nil {
		panic(err)
	}
	userId := bot.UserID

	return &SlackHandler{
		SigningSecret: config.SlackSigningSecret,
		client:        client,
		userId:        userId,
		service:       service,
	}
}

func (h *SlackHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/command", h.HandleCommands)
	e.POST("/actions", h.HandleActions)
}

func (h *SlackHandler) HandleCommands(c echo.Context) error {
	header := c.Request().Header

	verifier, err := slack.NewSecretsVerifier(header, h.SigningSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Request().Body = ioutil.NopCloser(io.TeeReader(c.Request().Body, &verifier))
	slackCommand, err := slack.SlashCommandParse(c.Request())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err = verifier.Ensure(); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if slackCommand.Command != "/도서관" {
		return echo.NewHTTPError(http.StatusBadRequest, "Command not supported")
	}

	userName := slackCommand.UserName
	command := strings.Split(slackCommand.Text, " ")
	switch command[0] {
	case "검색":
		if len(command) != 2 {
			return c.String(http.StatusOK, "명령이 잘못되었어요.\n사용 방법: /도서관 검색 `<검색어>`")
		}
		query := command[1]
		books, err := h.service.Search(query)
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		msg := views.RenderSearchResult(query, books)
		b, err := json.MarshalIndent(msg, "", "    ")
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		return c.JSONBlob(http.StatusOK, b)
	case "대출":
		if len(command) == 1 {
			return c.String(http.StatusOK, "명령이 잘못되었어요.\n사용 방법: /도서관 대출 `<책 이름의 일부>`")
		}

		query := strings.Join(command[1:], " ")
		books, err := h.service.Search(query)
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		if len(books) == 0 {
			return c.String(http.StatusOK, "해당 검색어로 책을 찾을 수 없었어요. 정확히 입력하셨는지 확인해주세요.")
		}

		if len(books) != 1 {
			return c.String(http.StatusOK, "해당 검색어로 찾은 책이 한 권이 아니에요. 제목을 더 자세하게 입력해주세요.")
		}

		book := books[0]
		book, err = h.service.Borrow(book, userName)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		msg := views.RenderBorrowResult(book)
		b, err := json.MarshalIndent(msg, "", "    ")
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		return c.JSONBlob(http.StatusOK, b)
	case "반납":
		if len(command) == 1 {
			return c.String(http.StatusOK, "명령이 잘못되었어요.\n사용 방법: /도서관 반납 `<책 이름의 일부>`")
		}

		query := strings.Join(command[1:], " ")
		books, err := h.service.Search(query)
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		if len(books) == 0 {
			return c.String(http.StatusOK, "해당 검색어로 책을 찾을 수 없었어요. 정확히 입력하셨는지 확인해주세요.")
		}

		if len(books) != 1 {
			return c.String(http.StatusOK, "해당 검색어로 찾은 책이 한 권이 아니에요. 제목을 더 자세하게 입력해주세요.")
		}

		book := books[0]
		book, err = h.service.Return(book, userName)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		msg := views.RenderReturnResult(book)
		b, err := json.MarshalIndent(msg, "", "    ")
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		return c.JSONBlob(http.StatusOK, b)
	case "연장":
		if len(command) == 1 {
			return c.String(http.StatusOK, "명령이 잘못되었어요.\n사용 방법: /도서관 연장 `<책 이름의 일부>`")
		}

		query := strings.Join(command[1:], " ")
		books, err := h.service.Search(query)
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		if len(books) == 0 {
			return c.String(http.StatusOK, "해당 검색어로 책을 찾을 수 없었어요. 정확히 입력하셨는지 확인해주세요.")
		}

		if len(books) != 1 {
			return c.String(http.StatusOK, "해당 검색어로 찾은 책이 한 권이 아니에요. 제목을 더 자세하게 입력해주세요.")
		}

		book := books[0]
		book, err = h.service.Extend(book, userName)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}

		msg := views.RenderExtendResult(book)
		b, err := json.MarshalIndent(msg, "", "    ")
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		return c.JSONBlob(http.StatusOK, b)
	case "현황":
		if len(command) != 1 {
			return c.String(http.StatusOK, "명령이 잘못되었어요.\n사용 방법: /도서관 현황")
		}

		books, err := h.service.Status(userName)
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		msg := views.RenderStatusResult(books, userName)
		b, err := json.MarshalIndent(msg, "", "    ")
		if err != nil {
			return c.String(http.StatusOK, "서버 오류가 발생했어요. :( 나중에 다시 시도하세요.")
		}

		return c.JSONBlob(http.StatusOK, b)
	default:
		return c.String(http.StatusOK, "잘못된 명령어예요. `검색`, `대출`, `반납`, `연장`, `현황` 중 하나를 선택해주세요.")
	}
}

func (h *SlackHandler) HandleActions(c echo.Context) error {
	return nil
}
