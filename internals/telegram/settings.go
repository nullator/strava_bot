package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleSettingsComand(message *tgbotapi.Message) error {
	id := message.Chat.ID

	msg_txt := "Для загрузки файла тренировки в страву необходимо указывать " +
		"ряд параметров: тип тренировки (велосипед, спорт, лыжи и т.п.), название, " +
		"описание и другие.\n\n" +
		"Имеется возможность загружать тренировки в страву автоматически при получении " +
		"файла используя настройки по умолчанию (например, каждый раз указывать " +
		"тип тренировки \"Велозаезд\"), или перед загрузкой файла в страву вручную " +
		"вводить данные тренировки.\n\n" +
		"Выбери способ загрузки файла транировки (настройку можно поменять позже в любое время):\n"

	msg := tgbotapi.NewMessage(id, msg_txt)
	msg.ParseMode = "Markdown"

	var keyboard = tgbotapi.NewInlineKeyboardMarkup()
	var buttons = tgbotapi.NewInlineKeyboardRow()

	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Default", "use_def"))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Variables", "use_var"))
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, buttons)

	msg.ReplyMarkup = keyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
