package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
	"github.com/qrlzvrn/Clozapinum/erro"
)

//CallbackHandler - перехватывает сообщения от нажатий на inlineKeyboard и выдает один или несколько конфигов ответных сообщений
func CallbackHandler(callbackQuery *tgbotapi.CallbackQuery) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, erro.Err) {
	tguserID := callbackQuery.From.ID
	conn, err := db.ConnectToBD()
	if err != nil {
		return nil, nil, nil, err
	}
	defer conn.Close()

	switch callbackQuery.Data {
	case "createTask":
		msg, newKeyboard, newText, err := CreateTask(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "createCategory":
		msg, newKeyboard, newText, err := CreateCategory(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "backToAllCategoriesKeyboard":
		msg, newKeyboard, newText, err := BackToAllCategoriesKeyboard(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "backToListTasks":
		msg, newKeyboard, newText, err := BackToListTasks(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "backToTask":
		msg, newKeyboard, newText, err := BackToTask(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "choose":
		msg, newKeyboard, newText, err := Choose(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "complete":
		msg, newKeyboard, newText, err := Complete(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "deleteTask":
		msg, newKeyboard, newText, err := DeleteTask(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "deleteCategoryQuestion":
		msg, newKeyboard, newText := ConfirmDeleteCategory(callbackQuery, conn, tguserID)

		return msg, newKeyboard, newText, nil
	case "deleteCategory":
		msg, newKeyboard, newText, err := DeleteCategory(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "change":
		msg, newKeyboard, newText, err := ChangeTask(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "changeTitle":
		msg, newKeyboard, newText, err := ChangeTitle(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "changeDescription":
		msg, newKeyboard, newText, err := ChangeDescription(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	case "changeDeadline":
		msg, newKeyboard, newText, err := ChangeDeadline(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	default:
		msg, newKeyboard, newText, err := ViewCategory(callbackQuery, conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		return msg, newKeyboard, newText, nil
	}
}
