package telegram

import (
	"fmt"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotService) douwnloadFile(document *tgbotapi.Document) (string, string, error) {
	// check file type
	file_name := document.FileName
	slog.Info("get file from user",
		slog.String("file_name", file_name))
	l := len([]rune(file_name))

	file_format := string([]rune(file_name)[l-4 : l])
	file_format = strings.ToLower(file_format)
	err_comment := fmt.Sprintf("Получен файл некорректного формата (%s), "+
		"принимаются файлы только формата .txt, .gpx и .fit", file_format)
	if file_format != ".fit" && file_format != ".tcx" && file_format != ".gpx" {
		return "", err_comment, fmt.Errorf("incorrect file format (%s)", file_format)
	}

	//download file
	err := b.service.Telegram.GetFile(document.FileName, document.FileID)
	if err != nil {
		err_comment := "Произошла внутренняя ошибка сервера Telegram"
		return "", err_comment, fmt.Errorf("error download fit file (%s)", err.Error())
	}

	file_path := fmt.Sprintf("activity/%s", document.FileName)
	return file_path, "", nil
}
