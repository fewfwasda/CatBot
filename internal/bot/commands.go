package bot

import (
	"CatBot/internal/data"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// SetCommand  Регистрирует команды бота
func SetCommand(bot *tgbotapi.BotAPI) error {
	commands := []tgbotapi.BotCommand{
		{
			Command:     "about",
			Description: "Информация о кото-боте",
		},
		{
			Command:     "showcats",
			Description: "Показать котиков",
		},
	}
	setCommands := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(setCommands)
	if err != nil {
		log.Printf("Не удалось установить команды: %v", err)
		return err
	}
	return nil
}

// CreateInlineMenu формирует сообщение с inline-клавиатурой для выбора котиков.
func CreateInlineMenu(chatID int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, data.ShowCatsText)

	whiteCat := tgbotapi.NewInlineKeyboardButtonData("Беленький котик🐻‍❄️", data.WhiteCat)
	blackCat := tgbotapi.NewInlineKeyboardButtonData("Черненький кот🐈‍⬛", data.BlackCat)
	gingerCat := tgbotapi.NewInlineKeyboardButtonData("Рыжий бесстыжий котик🐯", data.GingerCat)
	randomCat := tgbotapi.NewInlineKeyboardButtonData("Случайный котик🎲", data.RandomCat)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(whiteCat, blackCat),
		tgbotapi.NewInlineKeyboardRow(gingerCat, randomCat),
	)
	msg.ReplyMarkup = keyboard
	return msg
}
