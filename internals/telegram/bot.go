package telegram

import (
	"fmt"
	"log/slog"
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
	b.service.Logger.Info("Authorized", slog.String("account", b.bot.Self.UserName))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	b.handleUpdates(updates)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			sender := fmt.Sprintf("%s (%s)",
				update.Message.From.UserName,
				update.Message.From.String())
			b.service.Logger.Info("Получено сообщение",
				slog.String("sender", sender),
				slog.String("message", update.Message.Text))
			if update.Message.IsCommand() {
				b.service.Logger.Info("Зафиксирована команда",
					slog.String("sender", sender),
					slog.String("cmd", update.Message.Text))
				if err := b.handleCommand(update.Message); err != nil {
					b.service.Logger.Error("При обработке команды произошла ошибка",
						slog.String("cmd", update.Message.Command()),
						slog.String("error", err.Error()),
					)
				}
				continue
			}

			d := update.Message.Document
			id := update.Message.From.ID
			// if get file (d) from user
			if d != nil {
				filename, err := b.douwnloadFile(d)
				if err != nil {
					b.service.Logger.Error("error download file",
						slog.String("filename", filename),
						slog.String("error", err.Error()),
					)
					msg := tgbotapi.NewMessage(id, "Не удалось обработать файл")
					msg.ParseMode = "Markdown"
					_, err := b.bot.Send(msg)
					if err != nil {
						b.service.Logger.Error("error send message to user",
							slog.String("error", err.Error()))
					}
				} else {
					// upload file to strava
					err = b.service.UploadActivity(filename, id)
					var msg_txt string
					if err != nil {
						b.service.Logger.Error("error upload file to Strava server:",
							slog.String("error", err.Error()))
						msg_txt = "Произошла ошибка загрузки файла на сервер Strava.\n" +
							"Проверьте корректность файла (поддерживаются файлы .fit, .tcx и .gpx)" +
							" или попробуйте повторить загрузку позже."
					} else {
						msg_txt = "Файл тренировки успешно обработан и загружен в Strava"
					}
					msg := tgbotapi.NewMessage(id, msg_txt)
					msg.ParseMode = "Markdown"
					_, err = b.bot.Send(msg)
					if err != nil {
						b.service.Logger.Error("error send message to user",
							slog.String("error", err.Error()))
					}
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
		b.service.Logger.Error("error send message to user:",
			slog.String("error", err.Error()))
	}

}
