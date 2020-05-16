package handlers

import (
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	db "github.com/qrlzvrn/Clozapinum/db"
	"github.com/qrlzvrn/Clozapinum/errorz"
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
		return nil, nil, nil, err
	}
	defer conn.Close()

	//Проверяем наличие пользователя в базе,
	//если пользователь существует, то переходим к проверке категорий
	//если нет, тогда добавляем его в базу и предлагаем создать новую категорию
	isExist, err := db.CheckUser(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isExist == false {
		err := db.CreateUser(conn, tguserID)
		if err != nil {
			return nil, nil, nil, err
		}

		msg = tgbotapi.NewMessage(message.Chat.ID, "Здравствуйте, для того, что бы начать работу, вам нужно создать вашу первую категорию. После этого вы сможете добавлять в нее новые задачи, а так же создавать новые категории и работать с ними. Для ознакомления с правилами и способами создания задач и категорий воспользуйтесь командой /help.\n\nА сейчас, давайте начнем и создадим новую категорию!\nВведите название:")
		newKeyboard = nil
		newText = nil

		err = db.ChangeUserState(conn, tguserID, "categoryCreation")
		if err != nil {
			return nil, nil, nil, err
		}
	}

	//Проверяем есть ли у пользователя категории
	//если категорий нет, то предлагаем создать,
	//если же они есть, то
	//выводим клавиатуру с категориями
	allCategories, err := db.ListAllCategories(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	if allCategories == nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Для того, что бы начать, нужно создать категорию. Введите название:")
		newKeyboard = nil
		newText = nil

		err = db.ChangeUserState(conn, tguserID, "categoryCreation")
		if err != nil {
			return nil, nil, nil, err
		}
	}

	allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)
	msgConf := tgbotapi.NewMessage(message.Chat.ID, "Ваши категории:")
	msgConf.ReplyMarkup = allCategoriesKeyboard
	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// Help ...
func Help(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {

	msg = tgbotapi.NewMessage(message.Chat.ID, "Небольшая справка:\n\nСоздание категории:\n\n--1 запуск--\n- Введите команду /start\n- Введите название вашей первой категории\n- Поздравляем, категория успешно создана!\n\n--Дальнейшее использование--\n- Вернитесь к своим категориям и нажмите на кнопку Добавить\n- Введите название вашей новой категории\n- Поздравляем, еще одна категория успешно создана\n\nСоздание задачи:\n\n( Необходимо наличие категории )\n- Переходим в категорию, в котрую вы хотите добавить новую задачу\n- Если категория пуста, то вам будет предложенно ввести вашу задачу\n- Если у вас уже есть задачи в категории, то просто нажмите на кнопку Добавить\n(( Правила и способы создания задач ))\n- Вы можете просто ввести текст задачи и отправить его боту\n- Вы можете добавить к вашей задаче еще и дедлайн, для этого, вам необходимо пропустить одну строку после заголовка вашей задачи ( Дважды нажмите на Enter или Ctr+Enter, если вы работаете за компьютером ) и ввести дату в формате дд.мм.гггг\n- Вы можете так же добавить не только дедлайн, но и описание, так же пропустив одну строкй\n- Вы можете добавить к задаче и дедлайн и описание вместе, или что-то одно.\n- Но если вы хотите добавить и дедлайн и описание, то учтите, что дедлайн должен всегда идти перед описанием.")
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText
}

// Default ...
func Default(message *tgbotapi.Message) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {
	msg = tgbotapi.NewMessage(message.Chat.ID, "Простите, я так не умею :с")
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText
}

// Функции идущие далее служат для обработки текстовых сообщений пользователей.
// Так как в messageHandler принятие решений строится на основании значения state,
// то название данных функций соответсвует пришедшему значению.

// CategoryCreationAct ...
func CategoryCreationAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	nameCategory := message.Text
	err := db.CreateCategory(conn, tguserID, nameCategory)
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}

	err = db.CreateTask(conn, categoryID, title)
	if err != nil {
		// Распаковываем ошибку
		e := err.(*errorz.Err)

		//Проверяем значение ошибки
		if e.Err.Error() == "TitleIsNotDeadline" || e.Err.Error() == "NilTitle" {
			msgConf := tgbotapi.NewMessage(message.Chat.ID, "Простите, но у задачи обязательно должно быть название. Попробуйте еще раз")
			msgConf.ReplyMarkup = keyboard.CreateTaskKeyboard

			msg = msgConf
			newKeyboard = nil
			newText = nil

			return msg, newKeyboard, newText, nil
		}

		if e.Err.Error() == "DateErr" {
			msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется, с вашим дедлайном что-то не так, вы точно ввели его в формате дд.мм.гггг? Может быть такой даты не существует? Попробуте еще раз")
			msgConf.ReplyMarkup = keyboard.CreateTaskKeyboard

			msg = msgConf
			newKeyboard = nil
			newText = nil

			return msg, newKeyboard, newText, nil
		}
		return msg, newKeyboard, newText, nil
	}

	allTasks, err := db.ListTasks(conn, categoryID)
	if err != nil {
		return nil, nil, nil, err
	}
	err = db.ChangeUserState(conn, tguserID, "taskSelection")
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}
	taskID := message.Text
	intTaskID, err := strconv.Atoi(taskID)
	if err != nil || intTaskID < 1 {
		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется вы ввели что-то не то. Вы должны ввести id задачи. Давайте попробуем еще раз.\n\nВведите id задачи:")
		msgConf.ReplyMarkup = keyboard.SelectTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
		return msg, newKeyboard, newText, nil
	}

	isExist, err := db.IsTaskExist(conn, categoryID, intTaskID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isExist == false {
		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется, у вас нет задачи с таким id. Давайте попробуем еще раз?\n\nВведите id задачи:")
		msgConf.ReplyMarkup = keyboard.SelectTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
	}

	text, realTaskID, err := db.ViewTask("select", conn, categoryID, intTaskID, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.ChangeSelectTask(conn, tguserID, realTaskID)
	if err != nil {
		return nil, nil, nil, err
	}

	msgConf := tgbotapi.NewMessage(message.Chat.ID, text)

	isComplete, err := db.IsComplete(conn, realTaskID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isComplete == false {
		msgConf.ReplyMarkup = keyboard.TaskKeyboard
	} else {
		msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
	}

	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// ChangedTaskTitleAct ...
func ChangedTaskTitleAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	taskID, err := db.CheckSelectTaskID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	newTitle := message.Text

	if len([]rune(newTitle)) > 255 {
		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется длина вашего заголовка слишком велика, попробуте сократить его")
		msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		return msg, newKeyboard, newText, nil
	}

	err = db.ChangeTaskTitle(conn, tguserID, taskID, newTitle)
	if err != nil {
		return nil, nil, nil, err
	}

	text, _, err := db.ViewTask("change", conn, categoryID, taskID, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	msgConf := tgbotapi.NewMessage(message.Chat.ID, "")

	isComplete, err := db.IsComplete(conn, taskID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isComplete == false {
		msgConf.ReplyMarkup = keyboard.TaskKeyboard
	} else {
		msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
	}

	msgConf.Text = "Заголовок вашей задачи успешно изменен!\n\n" + text
	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// ChangedTaskDescriptionAct ...
func ChangedTaskDescriptionAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	taskID, err := db.CheckSelectTaskID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	newDescription := message.Text

	err = db.ChangeTaskDescription(conn, tguserID, taskID, newDescription)
	if err != nil {
		return nil, nil, nil, err
	}

	text, _, err := db.ViewTask("change", conn, categoryID, taskID, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	msgConf := tgbotapi.NewMessage(message.Chat.ID, "")

	isComplete, err := db.IsComplete(conn, taskID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isComplete == false {
		msgConf.ReplyMarkup = keyboard.TaskKeyboard
	} else {
		msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
	}

	msgConf.Text = "Заголовок вашей задачи успешно изменен!\n\n" + text
	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// ChangedTaskDeadlineAct ...
func ChangedTaskDeadlineAct(message *tgbotapi.Message, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	conn, err := db.ConnectToBD()
	if err != nil {
		return nil, nil, nil, err
	}
	defer conn.Close()

	taskID, err := db.CheckSelectTaskID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	newDeadline := message.Text

	deadlineOK, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, newDeadline)
	if deadlineOK == false {
		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется ваш дедлайн не соответсвтует формату дд.мм.гггг. Попробуйте снова")
		msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil
		return msg, newKeyboard, newText, nil
	}

	layout := "02.01.2006"
	t, err := time.Parse(layout, newDeadline)
	if err != nil {
		msgConf := tgbotapi.NewMessage(message.Chat.ID, "Кажется такой даты не существует. Попросбуйте еще раз")
		msgConf.ReplyMarkup = keyboard.ChangeSomethingInTaskKeyboard

		msg = msgConf
		newKeyboard = nil
		newText = nil

		return msg, newKeyboard, newText, nil
	}

	fmtDeadline := t.Format("01-02-2006")

	err = db.ChangeTaskDeadline(conn, tguserID, taskID, fmtDeadline)
	if err != nil {
		return nil, nil, nil, err
	}

	text, _, err := db.ViewTask("change", conn, categoryID, taskID, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	msgConf := tgbotapi.NewMessage(message.Chat.ID, "")

	isComplete, err := db.IsComplete(conn, taskID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isComplete == false {
		msgConf.ReplyMarkup = keyboard.TaskKeyboard
	} else {
		msgConf.ReplyMarkup = keyboard.CompletedTaskKeyboard
	}

	msgConf.Text = "Дедлайн вашей задачи успешно изменен!\n\n" + text
	msg = msgConf
	newKeyboard = nil
	newText = nil

	return msg, newKeyboard, newText, nil
}

// Функции идущие далее служат для обработки нажатий пользователей на инлайн кнопки
// Все данные функции относятся к CallbackHandler

// CreateTask ...
func CreateTask(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	err := db.ChangeUserState(conn, tguserID, "taskCreation")
	if err != nil {
		return nil, nil, nil, err
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
		return nil, nil, nil, err
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
		return nil, nil, nil, err
	}
	allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allCategoriesKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Ваши категории:")

	err = db.ChangeUserState(conn, tguserID, "allCategories")
	if err != nil {
		return nil, nil, nil, err
	}

	return msg, newKeyboard, newText, nil
}

// BackToListTasks ...
func BackToListTasks(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {

	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	allTasks, err := db.ListTasks(conn, categoryID)
	if err != nil {
		return nil, nil, nil, err
	}

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.SelectedCategoryKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allTasks)

	err = db.ChangeUserState(conn, tguserID, "taskSelection")
	if err != nil {
		return nil, nil, nil, err
	}

	return msg, newKeyboard, newText, nil
}

// BackToTask ...
func BackToTask(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}
	taskID, err := db.CheckSelectTaskID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	text, _, err := db.ViewTask("change", conn, categoryID, taskID, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	isComplete, err := db.IsComplete(conn, taskID)
	if err != nil {
		return nil, nil, nil, err
	}

	if isComplete == false {
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.TaskKeyboard)
	}

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CompletedTaskKeyboard)

	msg = nil
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Задача успешно выполнена!\n\n"+text)

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
		return nil, nil, nil, err
	}
	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.CompleteTask(conn, taskID)
	if err != nil {
		return nil, nil, nil, err
	}

	text, _, err := db.ViewTask("change", conn, categoryID, taskID, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	isComplete, err := db.IsComplete(conn, taskID)
	if err == nil && isComplete == false {
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.TaskKeyboard)
	} else if err == nil && isComplete == true {
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CompletedTaskKeyboard)
	} else {
		return nil, nil, nil, err
	}
	msg = nil
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Задача успешно выполнена!\n\n"+text)

	return msg, newKeyboard, newText, nil
}

// DeleteTask ....
func DeleteTask(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	taskID, err := db.CheckSelectTaskID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.DeleteTask(conn, taskID, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	allTasks, err := db.ListTasks(conn, categoryID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.ChangeUserState(conn, tguserID, "taskSelection")
	if err != nil {
		return nil, nil, nil, err
	}

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.SelectedCategoryKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Задача успешно удалена!\n\n"+allTasks)

	return msg, newKeyboard, newText, nil
}

// ConfirmDeleteCategory ...
func ConfirmDeleteCategory(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable) {

	msg = nil
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.DeleteCategoryKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Вы уверены, что хотите удалить данную категорию?")

	return msg, newKeyboard, newText
}

// DeleteCategory ...
func DeleteCategory(callbackQuery *tgbotapi.CallbackQuery, conn *sqlx.DB, tguserID int) (tgbotapi.Chattable, tgbotapi.Chattable, tgbotapi.Chattable, error) {
	categoryID, err := db.CheckSelectCategoryID(conn, tguserID)
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.DeleteCategory(conn, tguserID, categoryID)
	if err != nil {
		return nil, nil, nil, err
	}

	allCategories, err := db.ListAllCategories(conn, tguserID)
	if allCategories == nil {
		msg = nil
		newKeyboard = nil
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Для того, что бы начать, нужно создать категорию. Введите название:")

		err = db.ChangeUserState(conn, tguserID, "categoryCreation")
		if err != nil {
			return nil, nil, nil, err
		}
	}

	msg = nil

	allCategoriesKeyboard := keyboard.CreateKeyboarWithAllCategories(allCategories)
	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allCategoriesKeyboard)
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Категория успешно удалена!\n\nВаши категории:")

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
		return nil, nil, nil, err
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
		return nil, nil, nil, err
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
		return nil, nil, nil, err
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
		err := errorz.NewErr("Atoi error")
		return nil, nil, nil, err
	}

	allTasks, err := db.ListTasks(conn, intCategoryID)
	if err != nil {
		return nil, nil, nil, err
	}

	if allTasks == "" {
		err = db.ChangeSelectCategory(conn, tguserID, intCategoryID)
		if err != nil {
			return nil, nil, nil, err
		}
		err = db.ChangeUserState(conn, tguserID, "taskCreation")
		if err != nil {
			return nil, nil, nil, err
		}
		msg = nil
		newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.CreateTaskKeyboard)
		newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, "Пока что категория пуста. Давайте добавим первую задачу.\n\nВведите название задачи:")
	}

	msg = nil

	newKeyboard = tgbotapi.NewEditMessageReplyMarkup(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, keyboard.SelectedCategoryKeyboard)

	allTasks, err = db.ListTasks(conn, intCategoryID)
	if err != nil {
		return nil, nil, nil, err
	}
	newText = tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, allTasks)

	err = db.ChangeUserState(conn, tguserID, "taskSelection")
	if err != nil {
		return nil, nil, nil, err
	}

	err = db.ChangeSelectCategory(conn, tguserID, intCategoryID)
	if err != nil {
		return nil, nil, nil, err
	}

	return msg, newKeyboard, newText, nil
}
