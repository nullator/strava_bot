package telegram

import (
	"fmt"
	"log"
	"strava_bot/internals/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
}

func NewBot(bot *tgbotapi.BotAPI, service *service.Service) *Bot {
	return &Bot{bot, service}
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s\n", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				log.Printf("[%s] ввёл команду %s",
					update.Message.From.UserName, update.Message.Text)
				if err := b.handleCommand(update.Message); err != nil {
					log.Printf("При обработке команды %s произошла ошибка %s",
						update.Message.Command(), err)
				}
				continue
			}

			d := update.Message.Document
			id := update.Message.From.ID
			if d != nil {
				filename, err := b.douwnloadFile(d)
				if err != nil {
					log.Println(err.Error())
					msg := tgbotapi.NewMessage(id, "Не удалось обработать файл")
					msg.ParseMode = "Markdown"
					_, err := b.bot.Send(msg)
					if err != nil {
						log.Println(err.Error())
					}
				} else {
					// upload file to strava
					err = b.service.UploadActivity(filename, id)
					var msg_txt string
					if err != nil {
						log.Println(err.Error())
						msg_txt = "Произошла ошибка загрузки файла на сервер Strava.\n" +
							"Проверьте корректность файла (поддерживаются файлы .fit, .tcx и .gpx)" +
							" и попробуйте повторить загрузку позже."
					} else {
						msg_txt = "Файл тренировки успешно обработан и загружен в Strava"
					}
					msg := tgbotapi.NewMessage(id, msg_txt)
					msg.ParseMode = "Markdown"
					_, err = b.bot.Send(msg)
					if err != nil {
						log.Println(err.Error())
					}

				}

			}

		} else if update.CallbackQuery != nil {
			q := update.CallbackQuery.Data
			switch q {
			case "use_def":
				log.Println("def")
				buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
				msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons})
				_, err := b.bot.Send(msg)
				if err != nil {
					log.Println(err)
				}

			case "use_var":
				log.Println("var")
				buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
				msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons})
				_, err := b.bot.Send(msg)
				if err != nil {
					log.Println(err)
				}

			}

		}
	}
}

// TODO: прикреплять картинку с примером шаринга файла
func (b *Bot) SuccsesAuth(id int64, username string) {
	msg_txt := fmt.Sprintf("Успешная авторизация в аккаунт *%s*\n"+
		"Для загрузки тренировки отправляй боту файлы в формате .fit, .tcx или .gpx\nVPN не требуется", username)
	msg := tgbotapi.NewMessage(id, msg_txt)
	msg.ParseMode = "Markdown"
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Println(err.Error())
	}

}
