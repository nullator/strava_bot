package telegram

import (
	"fmt"
	"log/slog"
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
		"scope=activity:write,read_all&" +
		"state=%d"
)

func (b *BotService) handleCommand(message *tgbotapi.Message) error {
	const path = "internal.telegram.handleCMD.handleCommand"

	command := strings.ToLower(message.Command())
	switch command {

	case commandStart:
		err := b.handleStartComand(message)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}

	return nil
}

func (b *BotService) handleStartComand(message *tgbotapi.Message) error {
	const path = "internal.telegram.handleCMD.handleStartComand"

	URL := fmt.Sprintf(strava_auth_URL, os.Getenv("STRAVA_CLIENT_ID"),
		os.Getenv("STRAVA_REDIRECT_URL"), message.Chat.ID)

	msg_text := fmt.Sprintf("Для авторизации перейди по ссылке:\n[https://strava.com/](%s)", URL)
	msg := tgbotapi.NewMessage(message.Chat.ID, msg_text)

	msg.ParseMode = "Markdown"
	_, err := b.bot.Send(msg)
	if err != nil {
		slog.Error("error send message to user",
			slog.String("msg_text", msg_text),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("%s: %w", path, err)
	}

	sender := fmt.Sprintf("%s (%s)",
		message.From.UserName,
		message.From.String())
	slog.Info("Пользователю отправлена ссылка для авторизации",
		slog.String("sender", sender))
	return nil
}
