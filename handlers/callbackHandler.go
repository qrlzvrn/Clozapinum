package handlers

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
	keyboard "github.com/qrlzvrn/Clozapinum/keyboard"
)

//CallbackHandler - перехватывает сообщения от нажатий на inlineKeyboard и выдает один или несколько конфигов ответных сообщений
func CallbackHandler(callbackQuery *tgbotapi.CallbackQuery) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	var msg, newKeyboard, newText tgbotapi.Chattable

	switch callbackQuery.Data {
	case "createTask":
		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		tguserID := callbackQuery.From.ID

		err = db.ChangeUserState(conn, tguserID, "taskCreation")
		if err != nil {
			log.Panic(err)
		}

		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CreateTaskKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Введите задачу:")

	case "createCategory":
		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()
		tguserID := callbackQuery.From.ID

		err = db.ChangeUserState(conn, tguserID, "categoryCreation")
		if err != nil {
			log.Panic(err)
		}

		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CreateCategoryKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Как будет называться категория?")

	case "backToAllCategoriesKeyboard":
		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		tguserID := callbackQuery.From.ID

		allCategories, err := db.ListAllCategories(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}
		allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)

		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allCategoriesKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Ваши категории:")

		err = db.ChangeUserState(conn, tguserID, "allCategories")
		if err != nil {
			log.Panic(err)
		}

	case "backToListTasks":

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		tguserID := callbackQuery.From.ID

		categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}
		allTasks, err := db.ListTasks(conn, categoryID)
		if err != nil {
			log.Panic(err)
		}

		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.SelectedCategoryKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allTasks)

		err = db.ChangeUserState(conn, tguserID, "taskSelection")
		if err != nil {
			log.Panic(err)
		}
	case "backToTask":
		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		tguserID := callbackQuery.From.ID
		categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}
		taskID, err := db.CheckSelectTaskID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		text, err := db.ViewTask(conn, categoryID, taskID, tguserID)
		if err != nil {
			log.Panic(err)
		} else {

			isComplete, err := db.IsComplete(conn, taskID)
			if err == nil && isComplete == false {
				newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.TaskKeyboard)
			} else if err == nil && isComplete == true {
				newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CompletedTaskKeyboard)
			} else {
				log.Panic(err)
			}
			msg = nil
			newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Задача успешно выполнена!\n\n"+text)
		}
	case "choose":
		//проверяем выполнена ли задача и в зависимости от этого выдаем клавиатуру
		msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите id задачи:")
		msgConf.ReplyMarkup = keyboard.SelectTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
	case "complete":
		//делаем запрос к базе и отмечаем задачу выполненной,
		//если все ок, сообщаем об успехе,
		//если что-то пошло не так, то пишем пользователю сообщение об ошибке
		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		tguserID := callbackQuery.From.ID
		taskID, err := db.CheckSelectTaskID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}
		categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		err = db.CompleteTask(conn, taskID)
		if err != nil {
			log.Panic(err)
		}

		text, err := db.ViewTask(conn, categoryID, taskID, tguserID)
		if err != nil {
			log.Panic(err)
		} else {

			isComplete, err := db.IsComplete(conn, taskID)
			if err == nil && isComplete == false {
				newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.TaskKeyboard)
			} else if err == nil && isComplete == true {
				newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CompletedTaskKeyboard)
			} else {
				log.Panic(err)
			}
			msg = nil
			newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Задача успешно выполнена!\n\n"+text)
		}
	case "deleteTask":
		tguserID := callbackQuery.From.ID

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		taskID, err := db.CheckSelectTaskID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		err = db.DeleteTask(conn, taskID, tguserID)
		if err != nil {
			log.Panic(err)
		}

		allTasks, err := db.ListTasks(conn, categoryID)
		if err != nil {
			log.Panic(err)
		}

		err = db.ChangeUserState(conn, tguserID, "taskSelection")
		if err != nil {
			log.Panic(err)
		}

		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.SelectedCategoryKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Задача успешно удалена!\n\n"+allTasks)
	case "deleteCategoryQuestion":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.DeleteCategoryKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Вы уверены, что хотите удалить данную категорию?")
	case "deleteCategory":
		tguserID := callbackQuery.From.ID

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}

		err = db.DeleteCategory(conn, tguserID, categoryID)
		if err != nil {
			log.Panic(err)
		}

		allCategories, err := db.ListAllCategories(conn, tguserID)
		if err == nil && allCategories == nil {
			msg = nil
			newKeyboard = nil
			newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Для того, что бы начать, нужно создать категорию. Введите название:")

			err = db.ChangeUserState(conn, tguserID, "categoryCreation")
			if err != nil {
				log.Panic(err)
			}
		} else if err == nil && allCategories != nil {

			msg = nil

			allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)
			newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allCategoriesKeyboard)
			newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Категория успешно удалена!\n\nВаши категории:")
		} else {
			log.Panic(err)
		}

	case "change":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.ChangeTaskKeyboard)
		newText = nil
	case "changeTitle":
		tguserID := callbackQuery.From.ID

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		err = db.ChangeUserState(conn, tguserID, "changedTaskTitle")
		if err != nil {
			log.Panic(err)
		}
		msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите новый заголовок задачи:")
		msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
	case "changeDescription":
		tguserID := callbackQuery.From.ID

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		err = db.ChangeUserState(conn, tguserID, "changedTaskDescription")
		if err != nil {
			log.Panic(err)
		}
		msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите новое описание:")
		msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
	case "changeDeadline":
		tguserID := callbackQuery.From.ID

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		err = db.ChangeUserState(conn, tguserID, "changedTaskDeadline")
		if err != nil {
			log.Panic(err)
		}
		msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите новый дедлайн:")
		msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
	default:
		categoryID := callbackQuery.Data
		tguserID := callbackQuery.From.ID
		intCategoryID, err := strconv.Atoi(categoryID)
		if err != nil {
			log.Panic(err)
		}

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		allTasks, err := db.ListTasks(conn, intCategoryID)
		if err == nil && allTasks == "" {
			err = db.ChangeSelectCategory(conn, tguserID, intCategoryID)
			if err != nil {
				log.Panic(err)
			}
			err = db.ChangeUserState(conn, tguserID, "taskCreation")
			if err != nil {
				log.Panic(err)
			}
			msg = nil
			newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CreateTaskKeyboard)
			newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Пока что категория пуста. Давайте добавим первую задачу.\n\nВведите название задачи:")
		} else if err != nil {
			log.Panic(err)
		} else {

			msg = nil

			newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.SelectedCategoryKeyboard)

			allTasks, err := db.ListTasks(conn, intCategoryID)
			if err != nil {
				log.Panic(err)
			}
			newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allTasks)

			err = db.ChangeUserState(conn, tguserID, "taskSelection")
			if err != nil {
				log.Panic(err)
			}

			err = db.ChangeSelectCategory(conn, tguserID, intCategoryID)
			if err != nil {
				log.Panic(err)
			}
		}

	}
	return msg, newKeyboard, newText
}
