package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MainKeyboard - Клавиатура главного меню, расположеная внизу
var MainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Просмотреть"),
	),
)

//AddKeyboard - Клавиатура предлагающая выбор пользователю, создать задачу или категорию
//появляется после нажатия на кнопку Добавить
var AddKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Задачу", "createTask"),
		tgbotapi.NewInlineKeyboardButtonData("Категорию", "createKategory"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "backToMainMenu"),
	),
)

//CreateKeyboarWithAllCategory - генерирует клавиатуру со всеми переданными категориями
//появляется после нажатия на Просмотреть или на связку Добавить + Задачу
func CreateKeyboarWithAllCategory(categorys []string) tgbotapi.InlineKeyboardMarkup {
	AllCategoryKeyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, category := range categorys {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(category, category)
		row = append(row, btn)
		AllCategoryKeyboard.InlineKeyboard = append(AllCategoryKeyboard.InlineKeyboard, row)
	}
	return AllCategoryKeyboard
}

// CreateTaskKeyboard - клавиатура отмены создания задачи
// появляется после выбора категории в которой польователь хочет создать задачу
// и предложения ввести название новой задачи
var CreateTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "backToAllCategoryKeyboard"),
	),
)

// CreateCategoryKeyboard - клавиатура отмены создания категории
// появляется после комбинации Создать + Категория
// когда пользователю предлагается ввести название новой категории
var CreateCategoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "backToAddKeyboard"),
	),
)

//ChooseTaskKeyboard - клавиатура позволяющая либо выбрать задчу, либо вернуться назад
//появляется после нажатия на Просмотреть
var ChooseTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Выбрать", "choose"),
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToListAllCategory"),
	),
)

//TaskKeyboard - клавиатура с действиями над задачей
//появляется после нажатия на Просмотреть и введения id задачи
var TaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Выполнить", "complete"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "delete"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить", "change"),
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToListTasks"),
	),
)

//ChangeTaskKeyboard - клавиатура редактирования задачи
//появляется после комбинации Просмотреть + введения id + Измнить
var ChangeTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Заголовок", "changeTitle"),
		tgbotapi.NewInlineKeyboardButtonData("Описание", "changeDescription"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTask"),
	),
)
