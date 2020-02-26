package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
)

//MessageHandler - перехватывает простые текстовые сообщения и выдает один или несколько конфигов ответных сообщений
func MessageHandler(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {

	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "start":
			msg, newKeyboard, newText, err := Start(message)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText

		case "help":
			msg, newKeyboard, newText, err := Help(message)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		default:
			msg, newKeyboard, newText, err := Default(message)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		}
	} else {
		//Обрабатываем обычные текстовые сообщения
		//реакцию основываем на значении state в таблице tguser
		tguserID := message.From.ID

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		state, err := db.CheckUserState(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}
		switch state {
		case "categoryCreation":
			msg, newKeyboard, newText, err := CategoryCreationAct(message, conn, tguserID)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		case "taskCreation":
			msg, newKeyboard, newText, err := TaskCreationAct(message, conn, tguserID)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		case "taskSelection":
			msg, newKeyboard, newText, err := TaskSelectionAct(message, conn, tguserID)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		case "changedTaskTitle":
			msg, newKeyboard, newText, err := ChangedTaskTitleAct(message, conn, tguserID)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		case "changedTaskDescription":
			msg, newKeyboard, newText, err := ChangedTaskDescriptionAct(message, conn, tguserID)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		case "changedTaskDeadline":
			msg, newKeyboard, newText, err := ChangedTaskDeadlineAct(message, conn, tguserID)
			if err != nil {
				log.Panic(err)
			}

			return msg, newKeyboard, newText
		}
	}

	return msg, newKeyboard, newText
}
