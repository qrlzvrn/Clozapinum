package handlers

import (
	"log"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	db "github.com/qrlzvrn/Clozapinum/db"
	keyboard "github.com/qrlzvrn/Clozapinum/keyboard"
)

// Данная часть пакета содержит в себе функции выступающие посредниками между
// пакетом db и функциями messageHandler и callbackHandler,
// а так же призвана повыстить читаемсть и расширяемость в будущем

// Обявим переменные для конфигов сообщений в самом начале,
// что бы не приходилось повторять их объявление в каждой функции
var msg, newKeyboard, newText tgbotapi.Chattable

// Реализация функционала команд /start, /help и тех, что будут добавлены в будущем

// Start ...
func Start(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

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
	return msg, newKeyboard, newText, nil
}

// Help ...
func Help(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	msg = tgbotapi.NewMessage(message.Chat.ID, "Спешу на помощь, подождите немного, я сделаю все, что в моих силах")
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// Default ...
func Default(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	msg = tgbotapi.NewMessage(message.Chat.ID, "Простите, я так не умею :с")
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// Функции идущие далее служат для обработки текстовых сообщений пользователей.
// Так как в messageHandler принятие решений строится на основании значения state,
// то название данных функций соответсвует пришедшему значению.

// CategoryCreationAct ...
func CategoryCreationAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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

	return msg, newKeyboard, newText, nil
}

// TaskCreationAct ...
func TaskCreationAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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

		return msg, newKeyboard, newText, nil
	}

	if problem == "dateErr" {
		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется, с вашим дедлайном что-то не так, вы точно ввели его в формате дд.мм.гггг? Может быть такой даты не существует? Попробуте еще раз")
		msgConf.ReplyMarkup = keyboard.CreateTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		return msg, newKeyboard, newText, nil
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

	return msg, newKeyboard, newText, nil
}

// TaskSelectionAct ...
func TaskSelectionAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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
	return msg, newKeyboard, newText, nil
}

// ChangedTaskTitleAct ...
func ChangedTaskTitleAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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
	return msg, newKeyboard, newText, nil
}

// ChangedTaskDescriptionAct ...
func ChangedTaskDescriptionAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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
	return msg, newKeyboard, newText, nil
}

// ChangedTaskDeadlineAct ...
func ChangedTaskDeadlineAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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
	return msg, newKeyboard, newText, nil
}

// Функции идущие далее служат для обработки нажатий пользователей на инлайн кнопки
// Все данные функции относятся к CallbackHandler

// CreateTask ...
func CreateTask(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	err := db.ChangeUserState(conn, tguserID, "taskCreation")
	if err != nil {
		log.Panic(err)
	}

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CreateTaskKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Введите задачу:")

	return msg, newKeyboard, newText, nil
}

// CreateCategory ...
func CreateCategory(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	err := db.ChangeUserState(conn, tguserID, "categoryCreation")
	if err != nil {
		log.Panic(err)
	}

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CreateCategoryKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Как будет называться категория?")

	return msg, newKeyboard, newText, nil
}

// BackToAllCategoriesKeyboard ...
func BackToAllCategoriesKeyboard(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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

	return msg, newKeyboard, newText, nil
}

// BackToListTasks ...
func BackToListTasks(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

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

	return msg, newKeyboard, newText, nil
}

// BackToTask ...
func BackToTask(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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

	return msg, newKeyboard, newText, nil
}

// Choose ...
func Choose(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите id задачи:")
	msgConf.ReplyMarkup = keyboard.SelectTaskKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// Complete ....
func Complete(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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

	return msg, newKeyboard, newText, nil
}

// DeleteTask ....
func DeleteTask(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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

	return msg, newKeyboard, newText, nil
}

// ConfirmDeleteCategory ...
func ConfirmDeleteCategory(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.DeleteCategoryKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Вы уверены, что хотите удалить данную категорию?")

	return msg, newKeyboard, newText, nil
}

// DeleteCategory ...
func DeleteCategory(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
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
	return msg, newKeyboard, newText, nil
}

// ChangeTask ...
func ChangeTask(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.ChangeTaskKeyboard)
	newText = nil

	return msg, newKeyboard, newText, nil
}

// ChangeTitle ...
func ChangeTitle(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	err := db.ChangeUserState(conn, tguserID, "changedTaskTitle")
	if err != nil {
		log.Panic(err)
	}
	msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите новый заголовок задачи:")
	msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil
	return msg, newKeyboard, newText, nil
}

// ChangeDescription ...
func ChangeDescription(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	err := db.ChangeUserState(conn, tguserID, "changedTaskDescription")
	if err != nil {
		log.Panic(err)
	}
	msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите новое описание:")
	msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil
	return msg, newKeyboard, newText, nil
}

// ChangeDeadline ...
func ChangeDeadline(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	err := db.ChangeUserState(conn, tguserID, "changedTaskDeadline")
	if err != nil {
		log.Panic(err)
	}
	msgConf := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "Введите новый дедлайн:")
	msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

	msg = msgConf
	newKeyboard = nil
	newText = nil
	return msg, newKeyboard, newText, nil
}

// ViewCategory ...
func ViewCategory(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	categoryID := callbackQuery.Data
	intCategoryID, err := strconv.Atoi(categoryID)
	if err != nil {
		log.Panic(err)
	}

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

	return msg, newKeyboard, newText, nil
}