package handlers

import (
	"log"
	"strconv"

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
			isExist, err := db.CheckUser(conn, tguserID)
			if err == nil && isExist == false {
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

			} else if err == nil && isExist == true {
				//Проверяем есть ли у пользователя категории
				//если категорий нет, то предлагаем создать,
				//если же они есть, то
				//выводим клавиатуру с категориями
				allCategories, err := db.ListAllCategories(conn, tguserID)
				if err == nil && allCategories == nil {
					msg = tgbotapi.NewMessage(message.Chat.ID, "Для того, что бы начать, нужно создать категорию. Введите название:")
					newKeyboard = nil
					newText = nil

					err = db.ChangeUserState(conn, tguserID, "categoryCreation")
					if err != nil {
						log.Panic(err)
					}
				} else if err == nil && allCategories != nil {

					allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)
					msgConf := tgbotapi.NewMessage(message.Chat.ID, "Ваши категории:")
					msgConf.ReplyMarkup = allCategoriesKeyboard
					msg = msgConf
					newKeyboard = nil
					newText = nil
				} else {
					log.Panic(err)
				}
			}

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
			problem, err := db.CreateTask(conn, categoryID, title)
			if err != nil {
				log.Panic(err)
			}

			if problem == "titleIsNotDeadline" || problem == "nilTitle" {
				msgConf := tgbotapi.NewMessage(message.Chat.ID, "Простите, но у задачи обязательно должно быть название. Попробуйте еще раз")
				msgConf.ReplyMarkup = keyboard.CreateTaskKeyboard

				msg = msgConf
				newKeyboard = nil
				newText = nil
			}

			if problem == "dateErr" {
				msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется, с вашим дедлайном что-то не так, вы точно ввели его в формате дд.мм.гггг? Попробуте еще раз")
				msgConf.ReplyMarkup = keyboard.CreateTaskKeyboard

				msg = msgConf
				newKeyboard = nil
				newText = nil
			}

			allTasks, err := db.ListTasks(conn, categoryID)
			if err != nil {
				log.Panic(err)
			}
			err = db.ChangeUserState(conn, tguserID, "taskSelection")
			if err != nil {
				log.Panic(err)
			}

			msgConf := tgbotapi.NewMessage(message.Chat.ID, "Задача успешно создана!\n\n"+allTasks)
			msgConf.ReplyMarkup = keyboard.SelectedCategoryKeyboard

			msg = msgConf
			newKeyboard = nil
			newText = nil

		case "taskSelection":
			categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			taskID := message.Text
			intTaskID, err := strconv.Atoi(taskID)
			if err != nil {
				msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется вы ввели что-то не то. Вы должны ввести id задачи. Давайте попробуем еще раз.\n\nВведите id задачи:")
				msgConf.ReplyMarkup = keyboard.SelectTaskKeyboard

				msg = msgConf
				newKeyboard = nil
				newText = nil
			} else {

				isExist, err := db.IsTaskExist(conn, categoryID, intTaskID)
				if err != nil {
					log.Panic(err)
				} else if isExist == false {
					msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется, у вас нет задачи с таким id. Давайте попробуем еще раз?\n\nВведите id задачи:")
					msgConf.ReplyMarkup = keyboard.SelectTaskKeyboard

					msg = msgConf
					newKeyboard = nil
					newText = nil
				} else if isExist == true {

					text, err := db.ViewTask(conn, categoryID, intTaskID, tguserID)
					if err != nil && text == "taskSelectionErr" {
						msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется, у вас нет задачи с таким id. Давайте попробуем еще раз?\n\nВведите id задачи:")
						msgConf.ReplyMarkup = keyboard.SelectTaskKeyboard

						msg = msgConf
						newKeyboard = nil
						newText = nil
					} else {
						err = db.ChangeSelectTask(conn, tguserID, intTaskID)
						if err != nil {
							log.Panic(err)
						}

						msgConf := tgbotapi.NewMessage(message.Chat.ID, text)

						isComplete, err := db.IsComplete(conn, intTaskID)
						if err == nil && isComplete == false {
							msgConf.ReplyMarkup = keyboard.TaskKeyboard
						} else if err == nil && isComplete == true {
							msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
						} else {
							log.Panic(err)
						}

						msg = msgConf
						newKeyboard = nil
						newText = nil
					}
				}
			}
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

//InlineQueryHandler - перехватывает сообщения от нажатий на inlineKeyboard и выдает один или несколько конфигов ответных сообщений
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
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.ChangeTaskKeyboard)
		newText = nil
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
		//
	case "changeTitle":
		//
	case "changeDescription":
		//
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
