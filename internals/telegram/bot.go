package telegram

import (
	"fmt"
	"log/slog"
	"os"
	"strava_bot/internals/service"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotInterface interface {
	SuccsesAuth(id int64, username string) error
	Start()
}

var _ BotInterface = (*Bot)(nil)

type Bot struct {
	BotInterface
}

func NewBot(bot *tgbotapi.BotAPI, service *service.Service) *Bot {
	return &Bot{
		BotInterface: NewBotInterface(bot, service),
	}
}

type BotService struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
}

func NewBotInterface(bot *tgbotapi.BotAPI, service *service.Service) *BotService {
	return &BotService{bot, service}
}

func (b *BotService) Start() {
	slog.Info("Authorized", slog.String("account", b.bot.Self.UserName))
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	b.handleUpdates(updates)
}

func (b *BotService) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		switch {
		// Processing chat messages
		case update.Message != nil:
			sender := fmt.Sprintf("%s (%s)",
				update.Message.From.UserName,
				update.Message.From.String())

			// Checking for command input
			if update.Message.IsCommand() {
				slog.Info("Зафиксирована команда",
					slog.String("sender", sender),
					slog.String("cmd", update.Message.Text))
				if err := b.handleCommand(update.Message); err != nil {
					slog.Error("При обработке команды произошла ошибка",
						slog.String("cmd", update.Message.Command()),
						slog.String("error", err.Error()),
					)
				}
				continue
			}

			slog.Info("Получено сообщение",
				slog.String("sender", sender),
				slog.String("message", update.Message.Text))

			// Checking to receive feedback
			user_status, err := b.service.Telegram.GetStatus(update.Message.From.ID)
			if err != nil {
				slog.Error("Ошибка чтения из БД данных о статусе пользователя",
					slog.String("error", err.Error()))
			}
			if user_status == "feedback" {
				admin_id, err := strconv.ParseInt(os.Getenv("TG_adminID"), 10, 64)
				if err != nil {
					slog.Error("при выполнении авторизации не удалось распарсить ID в TelegramId",
						slog.String("error", err.Error()))
				}

				msg_text := fmt.Sprintf("Получена обратная связь от %s содержания: %s",
					update.Message.From.String(), update.Message.Text)
				msg := tgbotapi.NewMessage(admin_id, msg_text)
				_, err = b.bot.Send(msg)
				if err != nil {
					slog.Error("Не удалось отправить обратную связь",
						slog.String("error", err.Error()))
				}

				forvard_msg := tgbotapi.NewForward(admin_id,
					update.Message.Chat.ID,
					update.Message.MessageID)
				_, err = b.bot.Send(forvard_msg)
				if err != nil {
					slog.Error("Не удалось переслать сообщение",
						slog.String("error", err.Error()))
				}

				msg_text = "Спасибо, сообщение отправлено разработчику! " +
					"При необходимости можно повторно ввести команду /feedback " +
					"и отправить ещё одно сообщение, в том числе можно отправить файлы, скриншоты и т.п."
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, msg_text)
				_, err = b.bot.Send(msg)
				if err != nil {
					slog.Error("Не удалось отправить обратную связь",
						slog.String("error", err.Error()))
				}

				err = b.service.Telegram.SaveStatus(update.Message.Chat.ID, "null")
				if err != nil {
					slog.Error("Ошибка сохранения в БД данных о статусе пользователя",
						slog.String("error", err.Error()))

				}
				slog.Debug("Успешно изменён статус в БД")
				break
			}

			// Trying to process a training file
			d := update.Message.Document
			id := update.Message.From.ID
			// if get file (d) from user
			if d != nil {
				filename, err_comment, err := b.douwnloadFile(d)
				if err != nil {
					slog.Error("error download file",
						slog.String("filename", filename),
						slog.String("error", err.Error()),
					)
					err_comment := fmt.Sprintf("Не удалось обработать файл. %s", err_comment)
					msg := tgbotapi.NewMessage(id, err_comment)
					msg.ParseMode = "Markdown"
					_, err := b.bot.Send(msg)
					if err != nil {
						slog.Error("error send message to user",
							slog.String("error", err.Error()))
					}
				} else {
					// upload file to strava
					err = b.service.UploadActivity(filename, id)
					var msg_txt string
					if err != nil {
						slog.Error("error upload file to Strava server:",
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
						slog.Error("error send message to user",
							slog.String("error", err.Error()))
					}
				}
			}
		}
	}
}

// TODO: прикреплять картинку с примером шаринга файла
func (b *BotService) SuccsesAuth(id int64, username string) error {
	msg_txt := fmt.Sprintf("Успешная авторизация в аккаунт *%s*\n"+
		"Для загрузки тренировки отправляй боту файлы в формате .fit, .tcx или .gpx\nVPN не требуется", username)
	msg := tgbotapi.NewMessage(id, msg_txt)
	msg.ParseMode = "Markdown"
	_, err := b.bot.Send(msg)
	if err != nil {
		slog.Error("error send message to user:",
			slog.String("error", err.Error()))
	}
	return err
}
