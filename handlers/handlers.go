package handlers

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
	keyboard "github.com/qrlzvrn/Clozapinum/keyboard"
)

//MessageHandler - перехватывает простые текстовые сообщения и выдает конфиг ответного сообщения
func MessageHandler(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	var msg, newKeyboard, newText tgbotapi.Chattable

	if message.IsCommand() {
		cmd := message.Command()
		switch cmd {
		case "start":
			tguserID := message.From.ID

			conn, err := db.ConnectToBD()
			if err != nil {
				log.Panic(err)
			}
			defer conn.Close()

			//Проверяем наличие пользователя в базе,
			//если пользователь существует, то переходим к проверке категорий
			//если нет, тогда добавляем его в базу и предлагаем создать новую категорию
			err = db.CheckUser(conn, tguserID)
			if err != nil {
				err := db.CreateUser(conn, tguserID)
				if err != nil {
					log.Panic(err)
				}

				msg = tgbotapi.NewMessage(message.Chat.ID, "Для того, что бы начать, нужно создать категорию. Введите название:")
				newKeyboard = nil
				newText = nil

				err = db.ChangeUserState(conn, tguserID, "categoryCreation")
				if err != nil {
					log.Panic(err)
				}

			}
			//Проверяем есть ли у пользователя категории
			//если категорий нет, то предлагаем создать,
			//если же они есть, то выводим клавиатуру с категориями
			allCategories, err := db.ListAllCategories(conn, tguserID)
			if err != nil {
				msg = tgbotapi.NewMessage(message.Chat.ID, "Для того, что бы начать, нужно создать категорию. Введите название:")
				newKeyboard = nil
				newText = nil

				err = db.ChangeUserState(conn, tguserID, "categoryCreation")
				if err != nil {
					log.Panic(err)
				}
			}
			allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)
			msgConf := tgbotapi.NewMessage(message.Chat.ID, "Ваши категории:")
			msgConf.ReplyMarkup = allCategoriesKeyboard

		case "help":
			msg = tgbotapi.NewMessage(message.Chat.ID, "Спешу на помощь, подождите немного, я сделаю все, что в моих силах")
			newKeyboard = nil
			newText = nil
		default:
			msg = tgbotapi.NewMessage(message.Chat.ID, "Простите, я так не умею :с")
			newKeyboard = nil
			newText = nil
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
			nameCategory := message.Text
			err := db.CreateCategory(conn, tguserID, nameCategory)
			if err != nil {
				log.Panic(err)
			}
			allCategories, err := db.ListAllCategories(conn, tguserID)
			allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)

			msgConf := tgbotapi.NewMessage(message.Chat.ID, "Категория успешно создана!\n\nВаши категории:")
			msgConf.ReplyMarkup = allCategoriesKeyboard
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

			msg = tgbotapi.NewMessage(message.Chat.ID, "Задача успешно создана!")

		case "taskSelection":
			categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			taskID := message.Text
			intTaskID, err := strconv.Atoi(taskID)
			if err != nil {
				log.Panic(err)
			}

			text, err := db.ViewTask(conn, categoryID, intTaskID, tguserID)
			if err != nil {
				log.Panic(err)
			}
			err = db.ChangeSelectTask(conn, tguserID, intTaskID)
			if err != nil {
				log.Panic(err)
			}

			msgConf := tgbotapi.NewMessage(message.Chat.ID, text)
			msgConf.ReplyMarkup = keyboard.TaskKeyboard

			msg = msgConf
			newKeyboard = nil
			newText = nil
		case "changedTaskTitle":
			//
		case "changedTaskDescribe":
			//
		case "changedTaskDeadline":
			//
		}

		//обработка сообщений польователей о названии, описании, дедлайну задачи или категории
		//делаем запрос к бд и в зависимости от значения поля state решаем, что делать
	}

	return msg, newKeyboard, newText
}

//InlineQueryHandler - перехватывает сообщения от нажатий на inlineKeyboard и
//выдает один или несколько конфигов ответных сообщений
func InlineQueryHandler(callbackQuery *tgbotapi.CallbackQuery) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	var msg, newKeyboard, newText tgbotapi.Chattable

	switch callbackQuery.Data {
	case "createTask":
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
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Выберите категорию:")

		err = db.CheckUser(conn, tguserID)
		if err != nil {
			err = db.CreateUser(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
		}
		err = db.ChangeUserState(conn, tguserID, "taskCreation")
		if err != nil {
			log.Panic(err)
		}

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
			err = db.CreateUser(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
		}
		err = db.ChangeUserState(conn, tguserID, "categoryCreation")
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

		err = db.ChangeUserState(conn, tguserID, "taskSelection")
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
	default:
		categoryID := callbackQuery.Data
		tguserID := callbackQuery.From.ID
		strCategoryID, err := strconv.Atoi(categoryID)
		if err != nil {
			log.Panic(err)
		}

		conn, err := db.ConnectToBD()
		if err != nil {
			log.Panic(err)
		}
		defer conn.Close()

		msg = nil

		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.ChooseTaskKeyboard)

		tasksSlice := []string{}

		allTasks, err := db.ListTasks(conn, strCategoryID)
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

		err = db.ChangeUserState(conn, tguserID, "taskSelection")
		if err != nil {
			log.Panic(err)
		}

		err = db.ChangeSelectCategory(conn, tguserID, strCategoryID)
		if err != nil {
			log.Panic(err)
		}
	}
	return msg, newKeyboard, newText
}
