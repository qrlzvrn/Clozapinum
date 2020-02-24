package handlers

import (
	"log"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	db "github.com/qrlzvrn/Clozapinum/db"
	keyboard "github.com/qrlzvrn/Clozapinum/keyboard"
)

//MessageHandler - перехватывает простые текстовые сообщения и выдает один или несколько конфигов ответных сообщений
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
				msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется, с вашим дедлайном что-то не так, вы точно ввели его в формате дд.мм.гггг? Может быть такой даты не существует? Попробуте еще раз")
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
			conn, err := db.ConnectToBD()
			if err != nil {
				log.Panic(err)
			}
			defer conn.Close()

			tguserID := message.From.ID
			taskID, err := db.CheckSelectTaskID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			newTitle := message.Text

			if len([]rune(newTitle)) > 255 {
				msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется длина вашего заголовка слишком велика, попробуте сократить его")
				msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

				msg = msgConf
				newKeyboard = nil
				newText = nil
			} else {

				err = db.ChangeTaskTitle(conn, tguserID, taskID, newTitle)
				if err != nil {
					log.Panic(err)
				} else {

					text, err := db.ViewTask(conn, categoryID, taskID, tguserID)
					if err != nil {
						log.Panic(err)
					} else {
						msgConf := tgbotapi.NewMessage(message.Chat.ID, "")

						isComplete, err := db.IsComplete(conn, taskID)
						if err == nil && isComplete == false {
							msgConf.ReplyMarkup = keyboard.TaskKeyboard
						} else if err == nil && isComplete == true {
							msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
						} else {
							log.Panic(err)
						}
						msgConf.Text = "Заголовок вашей задачи успешно изменен!\n\n" + text
						msg = msgConf
						newKeyboard = nil
						newText = nil
					}

				}
			}
		case "changedTaskDescription":
			conn, err := db.ConnectToBD()
			if err != nil {
				log.Panic(err)
			}
			defer conn.Close()

			tguserID := message.From.ID
			taskID, err := db.CheckSelectTaskID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			newDescription := message.Text

			err = db.ChangeTaskDescription(conn, tguserID, taskID, newDescription)
			if err != nil {
				log.Panic(err)
			} else {
				text, err := db.ViewTask(conn, categoryID, taskID, tguserID)
				if err != nil {
					log.Panic(err)
				} else {
					msgConf := tgbotapi.NewMessage(message.Chat.ID, "")

					isComplete, err := db.IsComplete(conn, taskID)
					if err == nil && isComplete == false {
						msgConf.ReplyMarkup = keyboard.TaskKeyboard
					} else if err == nil && isComplete == true {
						msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
					} else {
						log.Panic(err)
					}
					msgConf.Text = "Заголовок вашей задачи успешно изменен!\n\n" + text
					msg = msgConf
					newKeyboard = nil
					newText = nil
				}

			}
		case "changedTaskDeadline":
			conn, err := db.ConnectToBD()
			if err != nil {
				log.Panic(err)
			}
			defer conn.Close()

			tguserID := message.From.ID
			taskID, err := db.CheckSelectTaskID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
			if err != nil {
				log.Panic(err)
			}
			newDeadline := message.Text

			deadlineOK, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, newDeadline)
			if deadlineOK == false {
				msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется ваш дедлайн не соответсвтует формату дд.мм.гггг. Попробуйте снова")
				msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

				msg = msgConf
				newKeyboard = nil
				newText = nil
			} else {

				layout := "02.01.2006"
				t, err := time.Parse(layout, newDeadline)
				if err != nil {
					msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется такой даты не существует. Попросбуйте еще раз")
					msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

					msg = msgConf
					newKeyboard = nil
					newText = nil
				} else {
					fmtDeadline := t.Format("01-02-2006")

					err = db.ChangeTaskDeadline(conn, tguserID, taskID, fmtDeadline)
					if err != nil {
						log.Panic(err)
					} else {
						text, err := db.ViewTask(conn, categoryID, taskID, tguserID)
						if err != nil {
							log.Panic(err)
						} else {
							msgConf := tgbotapi.NewMessage(message.Chat.ID, "")

							isComplete, err := db.IsComplete(conn, taskID)
							if err == nil && isComplete == false {
								msgConf.ReplyMarkup = keyboard.TaskKeyboard
							} else if err == nil && isComplete == true {
								msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
							} else {
								log.Panic(err)
							}
							msgConf.Text = "Дедлайн вашей задачи успешно изменен!\n\n" + text
							msg = msgConf
							newKeyboard = nil
							newText = nil
						}

					}
				}
			}
		}
		//обработка сообщений польователей о названии, описании, дедлайну задачи или категории
		//делаем запрос к бд и в зависимости от значения поля state решаем, что делать
	}

	return msg, newKeyboard, newText
}
