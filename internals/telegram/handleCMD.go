package telegram

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart    = "start"
	commandGet      = "get"
	commandSettings = "settings"

	strava_auth_URL = "https://www.strava.com/oauth/authorize?" +
		"client_id=%s&" +
		"redirect_uri=%s&" +
		"response_type=code&" +
		"approval_prompt=auto&" +
		"scope=activity:write,profile:read_all,activity:read,read&" +
		"state=%d"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	command := strings.ToLower(message.Command())
	switch command {

	case commandStart:
		err := b.handleStartComand(message)
		if err != nil {
			return err
		}

	case commandGet:
		err := b.handleGetComand(message)
		if err != nil {
			return err
		}

	case commandSettings:
		err := b.handleSettingsComand(message)
		if err != nil {
			return err
		}

	}

	return nil
}

func (b *Bot) handleStartComand(message *tgbotapi.Message) error {
	URL := fmt.Sprintf(strava_auth_URL, os.Getenv("STRAVA_CLIENT_ID"),
		os.Getenv("STRAVA_REDIRECT_URL"), message.Chat.ID)

	// [user mention](tg://user?id=12345)
	msg_text := fmt.Sprintf("Для авторизации перейди по ссылке:\n[https://strava.com/](%s)", URL)
	msg := tgbotapi.NewMessage(message.Chat.ID, msg_text)

	msg.ParseMode = "Markdown"
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Println("Выполнена команда Start")
	return nil
}

func (b *Bot) handleGetComand(message *tgbotapi.Message) error {
	id := message.Chat.ID
	var msg tgbotapi.MessageConfig
	msg.ChatID = id

	res, err := b.service.Strava.RefreshToken(id)
	if err != nil {
		msg.Text = err.Error()
	} else {
		msg.Text = res
	}

	_, err = b.bot.Send(msg)
	return err

}
