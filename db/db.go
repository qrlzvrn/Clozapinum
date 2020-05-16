package bd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qrlzvrn/Clozapinum/config"
	"github.com/qrlzvrn/Clozapinum/errorz"
)

func ConnectToBD() (*sqlx.DB, error) {
	dbConf, err := config.NewDBConf()
	if err != nil {
		return nil, err
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Name)

	db, err := sqlx.Connect("postgres", dbInfo)
	if err != nil {
		err := errorz.NewErr("failed connecting to database")
		return nil, err
	}

	return db, nil
}

func CreateUser(db *sqlx.DB, tguserID int) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, state) VALUES ($1, $2)", tguserID, "borned")
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to create a new user")
		return err
	}
	return nil
}

func CheckUser(db *sqlx.DB, tguserID int) (bool, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		err := errorz.NewErr("failed to verify user existence")
		return false, err
	}
	return isExist, nil
}

func ChangeUserState(db *sqlx.DB, tguserID int, state string) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET state=$1 where id=$2", state, tguserID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to change user state")
		return err
	}
	return nil
}

func ChangeSelectCategory(db *sqlx.DB, tguserID int, categoryID int) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_category=$1 where id=$2", categoryID, tguserID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to change the selected category")
		return err
	}
	return nil
}

func ChangeSelectTask(db *sqlx.DB, tguserID int, taskID int) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_task=$1 where id=$2", taskID, tguserID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to change the selected task")
		return err
	}
	return nil
}

func CheckUserState(db *sqlx.DB, tguserID int) (string, error) {
	var state string
	err := db.QueryRow("SELECT state FROM tguser WHERE id=$1", tguserID).Scan(&state)
	if err != nil {
		err := errorz.NewErr("failed to get user state")
		return "", err
	}
	return state, nil
}

func CheckSelectTaskID(db *sqlx.DB, tguserID int) (int, error) {
	var taskID int
	err := db.QueryRow("SELECT select_task FROM tguser WHERE id=$1", tguserID).Scan(&taskID)
	if err != nil {
		err := errorz.NewErr("failed to get id of selected task")
		return 0, err
	}
	return taskID, nil
}

func CheckSelectCategoryID(db *sqlx.DB, tguserID int) (int, error) {
	var categoryID int
	err := db.QueryRow("SELECT select_category FROM tguser WHERE id=$1", tguserID).Scan(&categoryID)
	if err != nil {
		err := errorz.NewErr("failed get id of selected task")
		return 0, err
	}
	return categoryID, nil
}

func CreateCategory(db *sqlx.DB, tguserID int, name string) error {
	var categoryID int
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO category (name) VALUES ($1)", name)
	err := tx.QueryRow("SELECT id FROM category where name=$1", name).Scan(&categoryID)
	if err != nil {
		err := errorz.NewErr("failed to get id of just created category")
		return err
	}
	tx.MustExec("INSERT INTO category_tguser (category_id, tguser_id) VALUES ($1, $2)", categoryID, tguserID)
	err = tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to create a new category")
		return err
	}
	return nil
}

func DeleteCategory(db *sqlx.DB, tguserID int, categoryID int) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_category=NULL WHERE id=$1", tguserID)
	tx.MustExec("UPDATE tguser SET select_task=NULL WHERE id=$1", tguserID)
	tx.MustExec("DELETE FROM task WHERE category_id=$1", categoryID)
	tx.MustExec("DELETE FROM category_tguser WHERE category_id=$1", categoryID)
	tx.MustExec("DELETE FROM category WHERE id=$1", categoryID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to delete category")
		return err
	}

	return nil
}

func ListAllCategories(db *sqlx.DB, tguserID int) ([][]string, error) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM category_tguser WHERE tguser_id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		err := errorz.NewErr("Error accessing the database")
		return nil, err
	}

	if isExist == false {
		return nil, nil
	}

	rows, err := db.Query("SELECT category.id, category.name FROM category LEFT JOIN category_tguser ON category_tguser.category_id=category.id WHERE category_tguser.tguser_id=$1", tguserID)
	if err != nil {
		err := errorz.NewErr("failed to get the list of categories")
		return nil, err
	}

	allCategories := [][]string{}

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			err := errorz.NewErr("failed to get the list of categories")
			return nil, err
		}

		allCategories = append(allCategories, []string{strconv.Itoa(id), name})
		defer rows.Close()
	}
	return allCategories, nil
}

func CreateTask(db *sqlx.DB, categoryID int, text string) error {

	sliceTaskText := strings.SplitN(text, "\n\n", 3)
	var fmtDeadline string
	//Проверяем сколько разных строк пришло
	len := len(sliceTaskText)
	if len == 3 {
		title := sliceTaskText[0]
		deadline := sliceTaskText[1]
		description := sliceTaskText[2]
		var fmtDeadline string
		//Проверяем не попал ли deadline в title
		titleErr, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, title)
		if titleErr == true {
			err := errorz.NewErr("TitleIsNotDeadline")
			return err
		}
		//Проверяем введен ли title
		if title == "" {
			err := errorz.NewErr("NilTitle")
			return err
		}

		//Проверяем правильность введения deadline
		deadlineOK, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, deadline)
		if deadlineOK == false {
			if deadline == "-" || deadline == "" {

				fmtDeadline = "01.01.1998"
				layout := "02.01.2006"
				t, err := time.Parse(layout, fmtDeadline)
				if err != nil {
					err := errorz.NewErr("DateErr")
					return err
				}
				fmtDeadline = t.Format("01-02-2006")
			} else {
				err := errorz.NewErr("DateErr")
				return err
			}

		} else {

			layout := "02.01.2006"
			t, err := time.Parse(layout, deadline)
			if err != nil {
				err := errorz.NewErr("DateErr")
				return err
			}
			fmtDeadline = t.Format("01-02-2006")
		}
		//Проверяем введено ли описание
		if description == "" {
			description = "-"
		}

		tx := db.MustBegin()
		tx.MustExec("INSERT INTO task(title, complete, category_id, description, deadline) VALUES ($1, $2, $3, $4, $5)", title, false, categoryID, description, fmtDeadline)
		err := tx.Commit()
		if err != nil {
			err := errorz.NewErr("failed to create a new task")
			return err
		}

	} else if len == 2 {
		var deadline string
		var description string

		title := sliceTaskText[0]
		something := sliceTaskText[1]

		//Проверяем не попал ли deadline в title
		titleErr, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, title)
		if titleErr == true {
			err := errorz.NewErr("TitleIsNotDeadline")
			return err
		}
		//Проверяем введен ли title
		if title == "" {
			err := errorz.NewErr("NilTitle")
			return err
		}

		isItDeadline, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, something)
		if isItDeadline == true {
			deadline = something
			description = "-"

			//форматируем deadline
			layout := "02.01.2006"
			t, err := time.Parse(layout, deadline)
			if err != nil {
				err := errorz.NewErr("DateErr")
				return err
			}
			fmtDeadline = t.Format("01-02-2006")
		} else if isItDeadline == false {
			//Если проверка по шаблону дедлайна не пройдена
			//проверяем на пустую строку или прочерк
			//если проверка не прошла,
			//тогда считаем, что полученная строка является описаниием
			if something == "" || something == "-" {

				deadline = "01.01.1998"

				layout := "02.01.2006"
				t, err := time.Parse(layout, fmtDeadline)
				if err != nil {
					err := errorz.NewErr("DateErr")
					return err
				}
				fmtDeadline = t.Format("01-02-2006")

				description = "-"
			}

			description = something
			deadline = "01.01.1998"

			layout := "02.01.2006"
			t, err := time.Parse(layout, deadline)
			if err != nil {
				err := errorz.NewErr("DateErr")
				return err
			}
			fmtDeadline = t.Format("01-02-2006")

		}
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO task(title, complete, category_id, description, deadline) VALUES ($1, $2, $3, $4, $5)", title, false, categoryID, description, fmtDeadline)
		err := tx.Commit()
		if err != nil {
			err := errorz.NewErr("failed to create a new task")
			return err
		}

	} else if len == 1 {
		title := sliceTaskText[0]
		var deadline string
		var description string
		//Проверяем не попал ли deadline в title
		titleErr, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, title)
		if titleErr == true {
			err := errorz.NewErr("TitleIsNotDeadline")
			return err
		}
		//Проверяем введен ли title
		if title == "" {
			err := errorz.NewErr("NilTitle")
			return err
		}

		deadline = "01.01.1998"

		layout := "02.01.2006"
		t, err := time.Parse(layout, deadline)
		if err != nil {
			err := errorz.NewErr("DateErr")
			return err
		}
		fmtDeadline = t.Format("01-02-2006")

		description = "-"
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO task(title, complete, category_id, description, deadline) VALUES ($1, $2, $3, $4, $5)", title, false, categoryID, description, fmtDeadline)
		err = tx.Commit()
		if err != nil {
			err := errorz.NewErr("failed to create a new task")
			return err
		}
	}
	return nil
}

func ViewTask(mode string, db *sqlx.DB, categoryID int, taskID int, tguserID int) (string, int, error) {

	var text string

	var title string
	var complete bool
	var description string
	var deadline string

	var realTaskID int

	if mode == "select" {
		taskID--
		tx := db.MustBegin()
		err := tx.QueryRow("SELECT id, title, description, complete, deadline-now()::date as deadline FROM task WHERE category_id=$1 ORDER BY complete LIMIT 1 OFFSET $2", categoryID, taskID).Scan(&realTaskID, &title, &description, &complete, &deadline)
		if err != nil {
			err := errorz.NewErr("failed to get task")
			return "", 0, err
		}
		tx.Commit()

	} else if mode == "change" {
		tx := db.MustBegin()
		err := tx.QueryRow("SELECT id, title, description, complete, deadline-now()::date as deadline FROM task WHERE category_id=$1 AND id=$2", categoryID, taskID).Scan(&realTaskID, &title, &description, &complete, &deadline)
		if err != nil {
			err := errorz.NewErr("failed to get task")
			return "", 0, err
		}
		tx.Commit()

	}

	nilDeadline, _ := regexp.MatchString(`^\-\d+`, deadline)

	if nilDeadline == true {
		if description == "-" {
			if complete == true {
				text = "\xF0\x9F\x93\x97 \t" + title
			} else {
				text = "\xF0\x9F\x93\x95 \t" + title
			}
		} else {
			if complete == true {
				text = "\xF0\x9F\x93\x97 \t" + title + "\n\n" + description
			} else {
				text = "\xF0\x9F\x93\x95 \t" + title + "\n\n" + description
			}
		}
	} else {

		if description == "-" {
			if complete == true {
				text = "\xF0\x9F\x93\x97 \t" + title + "\n\n" + "дедлайн: " + deadline + " дней"
			} else {
				text = "\xF0\x9F\x93\x95 \t" + title + "\n\n" + "дедлайн: " + deadline + " дней"
			}
		} else {
			if complete == true {
				text = "\xF0\x9F\x93\x97 \t" + title + "\n\n" + "дедлайн: " + deadline + " дней" + "\n\n" + description
			} else {
				text = "\xF0\x9F\x93\x95 \t" + title + "\n\n" + "дедлайн: " + deadline + " дней" + "\n\n" + description
			}
		}
	}

	return text, realTaskID, nil
}

func ChangeTaskTitle(db *sqlx.DB, tguserID int, taskID int, text string) error {
	tx := db.MustBegin()
	db.MustExec("UPDATE task SET title=$1 WHERE id=$2", text, taskID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to change task title")
		return err
	}
	return nil
}

func ChangeTaskDescription(db *sqlx.DB, tguserID int, taskID int, text string) error {
	tx := db.MustBegin()
	db.MustExec("UPDATE task SET description=$1 WHERE id=$2", text, taskID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to change task description")
		return err
	}
	return nil
}

func ChangeTaskDeadline(db *sqlx.DB, tguserID int, taskID int, text string) error {
	tx := db.MustBegin()
	db.MustExec("UPDATE task SET deadline=$1 WHERE id=$2", text, taskID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to change task deadline")
		return err
	}
	return nil
}

func ListTasks(db *sqlx.DB, categoryID int) (string, error) {
	var isExist string
	err := db.QueryRow("SELECT exists (SELECT 1 FROM task WHERE category_id=$1)", categoryID).Scan(&isExist)
	if err != nil {
		err := errorz.NewErr("failed to verify the existence of tasks in this category")
		return "", err
	}
	tasks := [][]string{}
	if isExist == "true" {
		rows, err := db.Query("SELECT title, complete FROM task WHERE category_id=$1 ORDER BY complete", categoryID)
		if err != nil {
			err := errorz.NewErr("failed to get task list")
			return "", err
		}

		id := 1

		for rows.Next() {
			var title string
			var complete bool

			rows.Scan(&title, &complete)
			tasks = append(tasks, []string{strconv.Itoa(id), title, strconv.FormatBool(complete)})
			id++
			defer rows.Close()
		}
	} else {
		return "", nil
	}

	tasksSlice := []string{}

	for _, task := range tasks {
		id := task[0]
		title := task[1]
		complete := task[2]
		if complete == "true" {
			text := "\xF0\x9F\x93\x97 ( " + id + " )\t" + title
			tasksSlice = append(tasksSlice, text)
		} else {
			text := "\xF0\x9F\x93\x95 ( " + id + " )\t" + title
			tasksSlice = append(tasksSlice, text)
		}
	}

	allTasksMsg := strings.Join(tasksSlice, "\n\n")

	return allTasksMsg, nil
}

func DeleteTask(db *sqlx.DB, taskID int, tguserID int) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_task=NULL WHERE id=$1", tguserID)
	tx.MustExec("DELETE FROM task WHERE id=$1", taskID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to delete the task")
		return err
	}
	return nil
}

func CompleteTask(db *sqlx.DB, taskID int) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE task SET complete=$1 WHERE id=$2", true, taskID)
	err := tx.Commit()
	if err != nil {
		err := errorz.NewErr("failed to complete the task")
		return err
	}
	return nil
}

func IsComplete(db *sqlx.DB, taskID int) (bool, error) {
	var isComplete bool
	tx := db.MustBegin()
	err := tx.QueryRow("SELECT complete FROM task where id=$1", taskID).Scan(&isComplete)
	if err != nil {
		err := errorz.NewErr("failed to verify if task completed")
		return false, err
	}
	tx.Commit()
	return isComplete, nil
}

func IsTaskExist(db *sqlx.DB, categoryID int, taskID int) (bool, error) {
	var isExist bool
	taskID--
	err := db.QueryRow("SELECT exists (SELECT 1 FROM task WHERE category_id=$1 ORDER BY complete LIMIT 1 OFFSET $2)", categoryID, taskID).Scan(&isExist)
	if err != nil {
		err := errorz.NewErr("failed to verify the existence of the task")
		return false, err
	}
	return isExist, nil
}
