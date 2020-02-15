package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//CreateKeyboarWithAllCategories - генерирует клавиатуру со всеми переданными категориями
//появляется после нажатия на Просмотреть или на связку Добавить + Задачу
func CreateKeyboarWithAllCategories(categories [][]string) tgbotapi.InlineKeyboardMarkup {
	AllCategoriesKeyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, category := range categories {
		id := category[0]
		name := category[1]
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(name, id)
		row = append(row, btn)
		AllCategoriesKeyboard.InlineKeyboard = append(AllCategoriesKeyboard.InlineKeyboard, row)
	}

	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData("Добавить", "createCategory")
	row = append(row, btn)
	AllCategoriesKeyboard.InlineKeyboard = append(AllCategoriesKeyboard.InlineKeyboard, row)
	return AllCategoriesKeyboard
}

// CreateTaskKeyboard - клавиатура отмены создания задачи
// появляется после выбора категории в которой польователь хочет создать задачу
// и предложения ввести название новой задачи
var CreateTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "backToAllCategoriesKeyboard"),
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

//SelectedCategoryKeyboard - клавиатура позволяющая либо выбрать задчу, либо вернуться назад
//появляется после нажатия на Просмотреть
var SelectedCategoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Выбрать", "choose"),
		tgbotapi.NewInlineKeyboardButtonData("Добавить", "createTask"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToAllCategoriesKeyboard"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить категорию", "deleteCategoryQuestion"),
	),
)

//DeleteCategoryKeyboard - клавиатура позволяющая согласиться с удалением категории
//или отказаться
var DeleteCategoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "deleteCategory"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "backToListTasks"),
	),
)

//SelectTaskKeyboard - клавиатура позволяющая отказаться от приглашения ввести id задачи
//после нажатия на кнопку выбрать в списке всех задач выбранной категории
var SelectTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToListTasks"),
	),
)

//TaskKeyboard - клавиатура с действиями над задачей
//появляется после нажатия на Просмотреть и введения id задачи
var TaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Выполнить", "complete"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "deleteTask"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить", "change"),
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToListTasks"),
	),
)

//CompletedTaskKeyboard - клавиатура с действиями над задачей
//появляется после нажатия на Просмотреть и введения id задачи,
//если задача уже выполненна
var CompletedTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить", "change"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "deleteTask"),
	),
	tgbotapi.NewInlineKeyboardRow(
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
