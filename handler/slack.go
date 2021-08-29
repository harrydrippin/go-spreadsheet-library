package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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
	e.POST("/action", h.HandleActions)
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
		if len(command) <= 1 {
			return c.String(http.StatusOK, "명령이 잘못되었어요.\n사용 방법: /도서관 검색 `<검색어>`")
		}

		query := strings.Join(command[1:], " ")
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
		return c.String(http.StatusOK, "잘못된 명령어예요. `검색`, `현황` 중 하나를 선택해주세요.")
	}
}

func (h *SlackHandler) HandleActions(c echo.Context) error {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(c.Request().FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v\n", err)
	}

	switch payload.Type {
	case slack.InteractionTypeBlockActions:
		for _, blockAction := range payload.ActionCallback.BlockActions {
			switch blockAction.ActionID {
			case utils.BorrowThisBook:
				book_id, err := strconv.Atoi(blockAction.Value)
				if err != nil {
					h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText("서버 오류가 발생했어요. :( 나중에 다시 시도하세요.", false))
					break
				}
				book, err := h.service.SearchById(book_id - 1)
				if err != nil {
					h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText("서버 오류가 발생했어요. :( 나중에 다시 시도하세요.", false))
					break
				}

				book, err = h.service.Borrow(book, payload.User.Name)
				if err != nil {
					_, err = h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText(err.Error(), false))
					break
				}

				msg := views.RenderBorrowResult(book)
				h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionBlocks(msg.Blocks.BlockSet...))

			case utils.ReturnThisBook:
				book_id, err := strconv.Atoi(blockAction.Value)
				if err != nil {
					h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText("서버 오류가 발생했어요. :( 나중에 다시 시도하세요.", false))
					break
				}
				book, err := h.service.SearchById(book_id - 1)
				if err != nil {
					h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText("서버 오류가 발생했어요. :( 나중에 다시 시도하세요.", false))
					break
				}

				book, err = h.service.Return(book, payload.User.Name)
				if err != nil {
					_, err = h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText(err.Error(), false))
					break
				}

				msg := views.RenderReturnResult(book)
				h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionBlocks(msg.Blocks.BlockSet...))
			case utils.ExtendThisBook:
				book_id, err := strconv.Atoi(blockAction.Value)
				if err != nil {
					h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText("서버 오류가 발생했어요. :( 나중에 다시 시도하세요.", false))
					break
				}
				book, err := h.service.SearchById(book_id - 1)
				if err != nil {
					h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText("서버 오류가 발생했어요. :( 나중에 다시 시도하세요.", false))
					break
				}

				book, err = h.service.Extend(book, payload.User.Name)
				if err != nil {
					_, err = h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionText(err.Error(), false))
					break
				}

				msg := views.RenderExtendResult(book)
				h.client.PostEphemeral(payload.Channel.ID, payload.User.ID, slack.MsgOptionBlocks(msg.Blocks.BlockSet...))
			}
		}
	}

	return c.String(http.StatusOK, "")
}
