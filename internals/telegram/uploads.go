package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) douwnloadFile(document *tgbotapi.Document) (string, error) {
	// check file type
	file_name := document.FileName
	l := len(file_name)
	file_format := string([]rune(file_name)[l-4 : l])
	if file_format != ".fit" && file_format != ".tcx" && file_format != ".gpx" {
		return "", fmt.Errorf("incorrect file format (%s)", file_format)
	}

	//download file
	err := b.service.Telegram.GetFile(document.FileName, document.FileID)
	if err != nil {
		return "", fmt.Errorf("error get file (%s)", err.Error())
	}

	file_path := fmt.Sprintf("activity/%s", document.FileName)
	return file_path, nil
}
