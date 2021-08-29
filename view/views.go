package view

import (
	"fmt"
	"strconv"

	models "github.com/harrydrippin/go-spreadsheet-library/model"
	"github.com/harrydrippin/go-spreadsheet-library/utils"
	"github.com/slack-go/slack"
)

func RenderSearchResult(query string, books []models.Book) slack.Message {
	spreadsheetLink := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", utils.NewConfig().GoogleSpreadsheetID)

	divSection := slack.NewDividerBlock()
	sections := make([]slack.Block, 0)

	// Header Text
	headerText := fmt.Sprintf("검색하신 *%s* 에 대한 *%d개* 의 결과가 있어요.", query, len(books))
	headerTextBlock := slack.NewTextBlockObject("mrkdwn", headerText, false, false)
	headerSection := slack.NewSectionBlock(headerTextBlock, nil, nil)
	sections = append(sections, headerSection)

	if len(books) == 0 {
		additionalText := fmt.Sprintf("조금 더 일반적인 키워드로 검색해보시거나, <%s|Spreadsheet>에서 직접 찾아주세요.", spreadsheetLink)
		additionalTextBlock := slack.NewTextBlockObject("mrkdwn", additionalText, false, false)
		additionalSection := slack.NewSectionBlock(additionalTextBlock, nil, nil)

		sections = append(sections, additionalSection)
		return slack.NewBlockMessage(sections...)
	} else if len(books) <= 5 {
		sections = append(sections, divSection)

		for _, book := range books {
			statusText := ""
			if book.Status == models.StatusBorrowed {
				statusText = fmt.Sprintf("*대출* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			} else if book.Status == models.StatusOverdue {
				statusText = fmt.Sprintf("*연체* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			} else {
				statusText = "사내 비치"
			}

			bookInfoText := fmt.Sprintf(">*%s*\n>%s 지음, %s\n>현재 상태: %s", book.Title, book.Author, book.Publisher, statusText)
			bookInfoBlock := slack.NewTextBlockObject("mrkdwn", bookInfoText, false, false)
			bookInfoSection := slack.NewSectionBlock(
				bookInfoBlock,
				nil,
				slack.NewAccessory(
					slack.NewButtonBlockElement(
						utils.BorrowThisBook,
						strconv.Itoa(book.ID),
						slack.NewTextBlockObject("plain_text", "대출하기", false, false),
					),
				),
			)
			sections = append(sections, bookInfoSection)
		}

		return slack.NewBlockMessage(sections...)
	} else {
		additionalText := fmt.Sprintf("검색 결과는 5개까지만 표시돼요. 그 이상은 <%s|Spreadsheet>에서 직접 찾아주세요.", spreadsheetLink)
		additionalTextBlock := slack.NewTextBlockObject("mrkdwn", additionalText, false, false)
		additionalSection := slack.NewSectionBlock(additionalTextBlock, nil, nil)
		sections = append(sections, additionalSection, divSection)

		for _, book := range books[:5] {
			statusText := ""
			if book.Status == models.StatusBorrowed {
				statusText = fmt.Sprintf("*대출* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			} else if book.Status == models.StatusOverdue {
				statusText = fmt.Sprintf("*연체* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			} else {
				statusText = "사내 비치"
			}

			bookInfoText := fmt.Sprintf(">*%s*\n>%s 지음, %s\n>현재 상태: %s", book.Title, book.Author, book.Publisher, statusText)
			bookInfoBlock := slack.NewTextBlockObject("mrkdwn", bookInfoText, false, false)
			bookInfoSection := slack.NewSectionBlock(
				bookInfoBlock,
				nil,
				slack.NewAccessory(
					slack.NewButtonBlockElement(
						utils.BorrowThisBook,
						strconv.Itoa(book.ID),
						slack.NewTextBlockObject("plain_text", "대출하기", false, false),
					),
				),
			)

			sections = append(sections, bookInfoSection)
		}

		return slack.NewBlockMessage(sections...)
	}
}

func RenderBorrowResult(book models.Book) slack.Message {
	divSection := slack.NewDividerBlock()
	sections := make([]slack.Block, 0)

	// Header Text
	headerText := "대출이 완료되었어요! 아래 내용을 확인해주세요.\n꼭 기한 내에 반납해주시고, 불가피하다면 연장 신청을 부탁드려요 :)"
	headerTextBlock := slack.NewTextBlockObject("mrkdwn", headerText, false, false)
	headerSection := slack.NewSectionBlock(headerTextBlock, nil, nil)
	sections = append(sections, headerSection, divSection)

	// Book Info
	statusText := fmt.Sprintf("*대출* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
	bookInfoText := fmt.Sprintf(">*%s*\n>%s 지음, %s\n>현재 상태: %s", book.Title, book.Author, book.Publisher, statusText)
	bookInfoBlock := slack.NewTextBlockObject("mrkdwn", bookInfoText, false, false)
	bookInfoSection := slack.NewSectionBlock(bookInfoBlock, nil, nil)

	sections = append(sections, bookInfoSection)

	return slack.NewBlockMessage(sections...)
}

func RenderReturnResult(book models.Book) slack.Message {
	divSection := slack.NewDividerBlock()
	sections := make([]slack.Block, 0)

	// Header Text
	headerText := "반납이 완료되었어요. 이용해주셔서 감사해요!"
	headerTextBlock := slack.NewTextBlockObject("mrkdwn", headerText, false, false)
	headerSection := slack.NewSectionBlock(headerTextBlock, nil, nil)
	sections = append(sections, headerSection, divSection)

	// Book Info
	statusText := "사내 비치"
	bookInfoText := fmt.Sprintf(">*%s*\n>%s 지음, %s\n>현재 상태: %s", book.Title, book.Author, book.Publisher, statusText)
	bookInfoBlock := slack.NewTextBlockObject("mrkdwn", bookInfoText, false, false)
	bookInfoSection := slack.NewSectionBlock(bookInfoBlock, nil, nil)

	sections = append(sections, bookInfoSection)

	return slack.NewBlockMessage(sections...)
}

func RenderExtendResult(book models.Book) slack.Message {
	divSection := slack.NewDividerBlock()
	sections := make([]slack.Block, 0)

	// Header Text
	headerText := "연장이 완료되었어요. 갱신된 기간을 확인해주시고, 기간 내에 반납해주세요!"
	headerTextBlock := slack.NewTextBlockObject("mrkdwn", headerText, false, false)
	headerSection := slack.NewSectionBlock(headerTextBlock, nil, nil)
	sections = append(sections, headerSection, divSection)

	// Book Info
	statusText := fmt.Sprintf("*연장 처리* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
	bookInfoText := fmt.Sprintf(">*%s*\n>%s 지음, %s\n>현재 상태: %s", book.Title, book.Author, book.Publisher, statusText)
	bookInfoBlock := slack.NewTextBlockObject("mrkdwn", bookInfoText, false, false)
	bookInfoSection := slack.NewSectionBlock(bookInfoBlock, nil, nil)

	sections = append(sections, bookInfoSection)

	return slack.NewBlockMessage(sections...)
}

func RenderStatusResult(books []models.Book, borrower string) slack.Message {
	spreadsheetLink := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", utils.NewConfig().GoogleSpreadsheetID)

	divSection := slack.NewDividerBlock()
	sections := make([]slack.Block, 0)

	// Header Text
	headerText := fmt.Sprintf("@%s 님께서는 %d 권의 책을 대출하셨어요.", borrower, len(books))
	headerTextBlock := slack.NewTextBlockObject("mrkdwn", headerText, false, false)
	headerSection := slack.NewSectionBlock(headerTextBlock, nil, nil)
	sections = append(sections, headerSection)

	if len(books) == 0 {
		additionalText := fmt.Sprintf("대출하시려면, `/도서관 검색 <책 이름의 일부>` 를 이용해보세요!")
		additionalTextBlock := slack.NewTextBlockObject("mrkdwn", additionalText, false, false)
		additionalSection := slack.NewSectionBlock(additionalTextBlock, nil, nil)

		sections = append(sections, additionalSection)
		return slack.NewBlockMessage(sections...)
	} else if len(books) <= 5 {
		sections = append(sections, divSection)

		for _, book := range books {
			statusText := ""
			if book.Status == models.StatusOverdue {
				statusText = fmt.Sprintf("*연체* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			} else {
				statusText = fmt.Sprintf("*대출* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			}

			bookInfoText := fmt.Sprintf(">*%s*\n>%s 지음, %s\n>현재 상태: %s", book.Title, book.Author, book.Publisher, statusText)
			bookInfoBlock := slack.NewTextBlockObject("mrkdwn", bookInfoText, false, false)
			bookInfoSection := slack.NewSectionBlock(bookInfoBlock, nil, nil)

			returnButtonBlock := slack.NewButtonBlockElement(
				utils.ReturnThisBook,
				strconv.Itoa(book.ID),
				slack.NewTextBlockObject("plain_text", "반납하기", false, false),
			)

			extendButtonBlock := slack.NewButtonBlockElement(
				utils.ExtendThisBook,
				strconv.Itoa(book.ID),
				slack.NewTextBlockObject("plain_text", "연장하기", false, false),
			)

			actionBlock := slack.NewActionBlock("status_action_block_"+strconv.Itoa(book.ID), returnButtonBlock, extendButtonBlock)

			sections = append(sections, bookInfoSection, actionBlock)
		}

		return slack.NewBlockMessage(sections...)
	} else {
		additionalText := fmt.Sprintf("검색 결과는 5개까지만 표시돼요. 그 이상은 <%s|Spreadsheet>에서 직접 찾아주세요.", spreadsheetLink)
		additionalTextBlock := slack.NewTextBlockObject("mrkdwn", additionalText, false, false)
		additionalSection := slack.NewSectionBlock(additionalTextBlock, nil, nil)
		sections = append(sections, additionalSection, divSection)

		for _, book := range books[:5] {
			statusText := ""
			if book.Status == models.StatusOverdue {
				statusText = fmt.Sprintf("*연체* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			} else {
				statusText = fmt.Sprintf("*대출* (%s, %s 반납 예정)", book.Borrower, book.DueDate)
			}

			bookInfoText := fmt.Sprintf(">*%s*\n>%s 지음, %s\n>현재 상태: %s", book.Title, book.Author, book.Publisher, statusText)
			bookInfoBlock := slack.NewTextBlockObject("mrkdwn", bookInfoText, false, false)
			bookInfoSection := slack.NewSectionBlock(bookInfoBlock, nil, nil)

			returnButtonBlock := slack.NewButtonBlockElement(
				utils.ReturnThisBook,
				strconv.Itoa(book.ID),
				slack.NewTextBlockObject("plain_text", "반납하기", false, false),
			)

			extendButtonBlock := slack.NewButtonBlockElement(
				utils.ExtendThisBook,
				strconv.Itoa(book.ID),
				slack.NewTextBlockObject("plain_text", "연장하기", false, false),
			)

			actionBlock := slack.NewActionBlock("status_action_block", returnButtonBlock, extendButtonBlock)

			sections = append(sections, bookInfoSection, actionBlock)
		}

		return slack.NewBlockMessage(sections...)
	}
}
