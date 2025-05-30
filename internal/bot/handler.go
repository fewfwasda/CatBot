package bot

import (
	"CatBot/internal/data"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

var lastCatMsgIDs = make(map[int64]int)
var catsPath = filepath.Join("..", "assets", "cats")

// HandleCallback Обработка нажатия inline кнопок
func HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, chatID int64) tgbotapi.CallbackConfig {

	if msgID, exists := lastCatMsgIDs[chatID]; exists {
		deleteMsg := tgbotapi.NewDeleteMessage(chatID, msgID)
		_, _ = bot.Request(deleteMsg)
	}

	var fileName string

	switch callback.Data {
	case data.WhiteCat, data.BlackCat, data.GingerCat:
		fileName = callback.Data + ".jpg"

	case data.RandomCat:
		files, err := os.ReadDir(catsPath)
		if err != nil || len(files) == 0 {
			log.Println("ошибка чтения каталога:", err)
			return tgbotapi.NewCallback(callback.ID, data.CatSentText)
		}

		fileName = files[rand.Intn(len(files))].Name()

	}

	if err := sendCatPhoto(bot, chatID, fileName); err != nil {
		log.Printf("ошибка отправки фото [%s]: %v", fileName, err)
	}

	callbackConfig := tgbotapi.NewCallback(callback.ID, data.CatSentText)
	return callbackConfig
}

// HandlerUserText Обработывает сообщение от пользователя и возвращает ответ
func HandlerUserText(update tgbotapi.Update) tgbotapi.MessageConfig {
	userText := update.Message.Text

	var answer tgbotapi.MessageConfig
	switch userText {

	case "/showcats":
		answer = CreateInlineMenu(update.Message.Chat.ID)
	case "/about":
		answer = tgbotapi.NewMessage(update.Message.Chat.ID, data.AboutText)
	default:
		answer = tgbotapi.NewMessage(update.Message.Chat.ID, data.WelcomeText)
	}
	return answer
}

// sendCatPhoto отправляет фото
func sendCatPhoto(bot *tgbotapi.BotAPI, chatID int64, filename string) error {
	fullPath := filepath.Join(catsPath, filename)
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(fullPath))
	if caption, ok := data.Captions[filename]; ok {
		photo.Caption = caption
	}
	sentMsg, err := bot.Send(photo)
	if err != nil {
		return err
	}
	lastCatMsgIDs[chatID] = sentMsg.MessageID
	return nil
}
