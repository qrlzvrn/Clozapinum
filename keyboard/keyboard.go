package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CreateKeyboarWithAllCategories - генерирует клавиатуру со всеми переданными категориями
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

// CreateTaskKeyboard - клавиатура создания задачи
var CreateTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "backToAllCategoriesKeyboard"),
	),
)

// CreateCategoryKeyboard - клавиатура создания категории
var CreateCategoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена", "backToAddKeyboard"),
	),
)

// SelectedCategoryKeyboard - клавиатура просмотра конкретной категории
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

// DeleteCategoryKeyboard - клавиатура позволяющая согласиться с удалением категории
// или отказаться
var DeleteCategoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Да", "deleteCategory"),
		tgbotapi.NewInlineKeyboardButtonData("Нет", "backToListTasks"),
	),
)

// SelectTaskKeyboard - клавиатура выбора задачи,
// позволяющая отказаться от приглашения ввести id задачи
var SelectTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToListTasks"),
	),
)

// TaskKeyboard - клавиатура с действиями над задачей
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

// CompletedTaskKeyboard - клавиатура с действиями над уже выполненной задачей
var CompletedTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить", "change"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "deleteTask"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToListTasks"),
	),
)

// ChangeTaskKeyboard - клавиатура, позволяющая выбрать,
// что конкретно изменить в данной задаче
var ChangeTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Заголовок", "changeTitle"),
		tgbotapi.NewInlineKeyboardButtonData("Описание", "changeDescription"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Дедлайн", "changeDeadline"),
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTask"),
	),
)

// ChangeSomethingInTaskKeyboard - клавиатура редактирования задачи,
// позволяет отменить редактирование задачи
var ChangeSomethingInTaskKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "backToTask"),
	),
)
