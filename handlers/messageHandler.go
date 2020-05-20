package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
)

//MessageHandler - перехватывает простые текстовые сообщения и выдает один или несколько конфигов ответных сообщений
func MessageHandler(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "start":
			msg, newKeyboard, newText, err := Start(message)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil

		case "help":
			msg, newKeyboard, newText := Help(message)

			return msg, newKeyboard, newText, nil
		default:
			msg, newKeyboard, newText := Default(message)

			return msg, newKeyboard, newText, nil
		}
	} else {
		//Обрабатываем обычные текстовые сообщения
		//реакцию основываем на значении state в таблице tguser
		tguserID := message.From.ID

		conn, err := db.ConnectToBD()
		if err != nil {
			return nil, nil, nil, err
		}
		defer conn.Close()

		state, err := db.CheckUserState(conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}
		switch state {
		case "categoryCreation":
			msg, newKeyboard, newText, err := CategoryCreationAct(message, conn, tguserID)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil
		case "taskCreation":
			msg, newKeyboard, newText, err := TaskCreationAct(message, conn, tguserID)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil
		case "taskSelection":
			msg, newKeyboard, newText, err := TaskSelectionAct(message, conn, tguserID)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil
		case "changedTaskTitle":
			msg, newKeyboard, newText, err := ChangedTaskTitleAct(message, conn, tguserID)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil
		case "changedTaskDescription":
			msg, newKeyboard, newText, err := ChangedTaskDescriptionAct(message, conn, tguserID)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil
		case "changedTaskDeadline":
			msg, newKeyboard, newText, err := ChangedTaskDeadlineAct(message, conn, tguserID)
			if err != nil {
				return nil, nil, nil, err
			}

			return msg, newKeyboard, newText, nil
		}
	}

	return msg, newKeyboard, newText, nil
}
