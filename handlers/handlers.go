package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/qrlzvrn/Clozapinum/db"
	"github.com/qrlzvrn/Clozapinum/keyboard"
)

//MessageHandler - перехватывает простые текстовые сообщения и выдает конфиг ответного сообщения
func MessageHandler(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	var msg, newKeyboard, newText tgbotapi.Chattable

	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "/start":
			msg = tgbotapi.NewMessage(message.Chat.ID, "Добрый вечер, это поленая информация. Можете нажать на кнопки, если хотите, конечно")
			newKeyboard = tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, message.MessageID, keyboard.MainKeyboard)
			newText = nil
		case "/help":
			msg = tgbotapi.NewMessage(message.Chat.ID, "Спешу на помощь, подождите немного, я сделаю все, что в моих силах")
			newKeyboard = nil
			newText = nil
		default:
			msg = tgbotapi.NewMessage(message.Chat.ID, "Простите, я так не умею :с")
			newKeyboard = nil
			newText = nil
		}
	} else {
		switch message.Text {
		case "создать":

		case "просмотреть":
			//
		case "Привет":
			id := message.From.ID

			conn, err := db.ConnectToBD()
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			db.CreateUser(conn, id)
			msg = tgbotapi.NewMessage(message.Chat.ID, "Все ОК")
			newKeyboard = nil
			newText = nil
		default:
			//обработка сообщений польователей о названии, описании, дедлайну задачи или категории
			//делаем запрос к бд и в зависимости от значения поля state решаем, что делать
		}
	}

	return msg, newKeyboard, newText
}

//InlineQueryHandler - перехватывает сообщения от нажатий на inlineKeyboard и
//выдает один или несколько конфигов ответных сообщений
func InlineQueryHandler(callbackQuery *tgbotapi.CallbackQuery) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	var msg, newKeyboard, newText tgbotapi.Chattable

	switch callbackQuery.Data {
	case "createTask":
		//
	case "createCategory":
		//
	case "backToMainMenu":
		//
	case "backToAllCategory":
		//
	case "backToAddKeyboard":
		//
	case "backToListAllCategory":
		//
	case "backToListTasks":
		//
	case "backToTask":
		//
	case "choose":
		//
	case "complete":
		//
	case "delete":
		//
	case "change":
		//
	case "changeTitle":
		//
	case "changeDescription":
		//
	}
	return msg, newKeyboard, newText
}
