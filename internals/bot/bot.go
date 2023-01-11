package bot

import (
	"strava_bot/pkg/base"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	base base.Base
}

func NewBot(bot *tgbotapi.BotAPI, base base.Base) *Bot {
	return &Bot{bot, base}
}
