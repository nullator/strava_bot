package telegram

import (
	"fmt"
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
		"scope=activity:write,read&" +
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
	}

	return nil
}

func (b *Bot) handleStartComand(message *tgbotapi.Message) error {
	URL := fmt.Sprintf(strava_auth_URL, os.Getenv("STRAVA_CLIENT_ID"),
		os.Getenv("STRAVA_REDIRECT_URL"), message.Chat.ID)

	msg_text := fmt.Sprintf("Для авторизации перейди по ссылке:\n[https://strava.com/](%s)", URL)
	msg := tgbotapi.NewMessage(message.Chat.ID, msg_text)

	msg.ParseMode = "Markdown"
	_, err := b.bot.Send(msg)
	if err != nil {
		b.service.Logger.Error("error send message to user: %v", err)
		return err
	}
	b.service.Logger.Info("Пользователю [%s (%s)] отправлены ссылка для авторизации",
		message.From.UserName, message.From.String())
	return nil
}
