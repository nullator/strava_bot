package telegram

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart     = "start"
	commandSubscribe = "subscribe"
	commandFeedback  = "feedback"

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
	msg := tgbotapi.NewMessage(message.Chat.ID, URL)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	log.Println("Выполнена команда Start")
	return nil
}
