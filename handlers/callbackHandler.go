package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
)

//CallbackHandler - перехватывает сообщения от нажатий на inlineKeyboard и выдает один или несколько конфигов ответных сообщений
func CallbackHandler(callbackQuery *tgbotapi.CallbackQuery) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	tguserID := callbackQuery.From.ID
	conn, err := db.ConnectToBD()
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	switch callbackQuery.Data {
	case "createTask":
		msg, newKeyboard, newText, err := CreateTask(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "createCategory":
		msg, newKeyboard, newText, err := CreateCategory(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "backToAllCategoriesKeyboard":
		msg, newKeyboard, newText, err := BackToAllCategoriesKeyboard(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "backToListTasks":
		msg, newKeyboard, newText, err := BackToListTasks(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "backToTask":
		msg, newKeyboard, newText, err := BackToTask(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "choose":
		msg, newKeyboard, newText, err := Choose(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "complete":
		msg, newKeyboard, newText, err := Complete(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "deleteTask":
		msg, newKeyboard, newText, err := DeleteTask(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "deleteCategoryQuestion":
		msg, newKeyboard, newText, err := ConfirmDeleteCategory(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "deleteCategory":
		msg, newKeyboard, newText, err := DeleteCategory(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "change":
		msg, newKeyboard, newText, err := ChangeTask(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "changeTitle":
		msg, newKeyboard, newText, err := ChangeTitle(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "changeDescription":
		msg, newKeyboard, newText, err := ChangeDescription(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	case "changeDeadline":
		msg, newKeyboard, newText, err := ChangeDeadline(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	default:
		msg, newKeyboard, newText, err := ViewCategory(callbackQuery, conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		return msg, newKeyboard, newText
	}
}
