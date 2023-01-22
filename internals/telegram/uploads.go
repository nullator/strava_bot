package telegram

import (
	"fmt"
	"strava_bot/internals/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) uploadActivity(document *tgbotapi.Document) (*models.UploadActivity, error) {
	var activity models.UploadActivity
	// check file type
	file_name := document.FileName
	l := len(file_name)
	file_format := string([]rune(file_name)[l-4 : l])
	switch file_format {
	case ".fit":
		activity.Data_type = file_format
	case ".tcx":
		activity.Data_type = file_format
	case ".gpx":
		activity.Data_type = file_format
	default:
		return nil, fmt.Errorf("incorrect file format (%s)", file_format)
	}

	//download file
	err := b.service.Telegram.GetFile(document.FileName, document.FileID)
	if err != nil {
		return nil, fmt.Errorf("error get file (%s)", err.Error())
	}

	activity.File = fmt.Sprintf("activity/%s", document.FileName)
	return &activity, nil
}
