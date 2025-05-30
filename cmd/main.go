package main

import (
	"CatBot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalf("Bot tolen not initialize: %v", botToken)
	}
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.SetCommand(botAPI)
	if err != nil {
		log.Fatalf("Не удалось установить команды: %v", err)
	}

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)

	for update := range updates {
		//Обработка inline стоки
		if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			callbackConfig := bot.HandleCallback(botAPI, update.CallbackQuery, chatID)
			_, err := botAPI.Request(callbackConfig)
			if err != nil {
				log.Printf("Ошибка Request: %v", err)
			}
		} else if update.Message != nil { //обработка обычных сообщений
			_, err := botAPI.Send(bot.HandlerUserText(update))
			if err != nil {
				log.Printf("Ошибка Send: %v", err)
			}
		}
	}
}
