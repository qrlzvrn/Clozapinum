package handlers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
	"github.com/qrlzvrn/Clozapinum/keyboard"
)

//MessageHandler - перехватывает простые текстовые сообщения и выдает конфиг ответного сообщения
func MessageHandler(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	var msg, newKeyboard, newText tgbotapi.Chattable

	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "/start":
			msgConf := tgbotapi.NewMessage(message.Chat.ID, "Добрый вечер, это поленая информация. Можете нажать на кнопки, если хотите, конечно")
			msgConf.ReplyMarkup = keyboard.MainKeyboard

			msg = msgConf
			newKeyboard = nil
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
			msgConf := tgbotapi.NewMessage(message.Chat.ID, "Что хотите создать?")
			msgConf.ReplyMarkup = keyboard.AddKeyboard

			msg = msgConf
			newKeyboard = nil
			newText = nil
		case "просмотреть":
			//
		case "Привет":
			id := message.From.ID

			conn, err := db.ConnectToBD()
			if err != nil {
				log.Panic(err)
			}
			defer conn.Close()

			db.CreateUser(conn, id)
			msg = tgbotapi.NewMessage(message.Chat.ID, "Все ОК")
			newKeyboard = nil
			newText = nil
		default:

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
				nameCategory := message.Text
				db.CreateCategory(conn, tguserID, nameCategory)
				msgConf := tgbotapi.NewMessage(message.Chat.ID, "Категория успешно создана")
				msgConf.ReplyMarkup = keyboard.MainKeyboard

				msg = msgConf
				newKeyboard = nil
				newText = nil
			case "taskCreation":
				tguserID := message.From.ID
				title := message.Text
				categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
				if err != nil {
					log.Panic(err)
				}
				err = db.CreateTask(conn, categoryID, title)
				if err != nil {
					log.Panic(err)
				}
			case "taskSelection":
				//
			case "taskChange":
				//
			}

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
	case "add":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.AddKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Что хотите создать?")
	case "ls":
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
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Что хотите создать?")

	case "createTask":
		//
	case "createCategory":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CreateCategoryKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Как будет называться категория?")

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()
		tguserID := callbackQuery.From.ID

		err = db.CheckUser(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}
		err = db.ChangeUserState(conn, tguserID, "categoryCreation")
		if err != nil {
			log.Panic(err)
		}

	case "backToMainMenu":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.MainKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Что будем делать сегодня?")

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()
		tguserID := callbackQuery.From.ID

		err = db.ChangeUserState(conn, tguserID, "mainMenu")
		if err != nil {
			log.Panic(err)
		}

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
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Ваши категории")

		err = db.ChangeUserState(conn, tguserID, "allCategories")
		if err != nil {
			log.Panic(err)
		}
	case "backToAddKeyboard":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.AddKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "что хотите создать?")

	case "backToListTasks":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.ChooseTaskKeyboard)

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		tasksSlice := []string{}

		tguserID := callbackQuery.From.ID

		categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
		if err != nil {
			log.Panic(err)
		}
		allTasks, err := db.ListTasks(conn, categoryID)
		if err != nil {
			log.Panic(err)
		}

		for _, task := range allTasks {
			id := task[0]
			title := task[1]
			complete := task[2]
			if complete == "true" {
				text := "\xE2\x9C\x85 #" + id + "-" + title
				tasksSlice = append(tasksSlice, text)
			} else {
				text := "\xE2\x9D\x8E #" + id + "-" + title
				tasksSlice = append(tasksSlice, text)
			}
		}

		allTasksMsg := strings.Join(tasksSlice, "\n")

		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allTasksMsg)

		err = db.ChangeUserState(conn, tguserID, "viewTaskList")
		if err != nil {
			log.Panic(err)
		}
	case "backToTask":
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.ChangeTaskKeyboard)
		newText = nil
	case "choose":
		//проверяем выполнена ли задача и в зависимости от этого выдаем клавиатуру
		msg = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите id задачи:")
		newKeyboard = nil
		newText = nil
	case "complete":
		//делаем запрос к базе и отмечаем задачу выполненной, если все ок, сообщаем об успехе, если что-то пошло не так, то пишем пользователю сообщение об ошибке
		msg = tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Задача выполненна")

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
