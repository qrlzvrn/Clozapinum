package bd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/qrlzvrn/Clozapinum/config"
	"github.com/qrlzvrn/Clozapinum/erro"
)

func ConnectToBD() (*sqlx.DB, erro.Err) {
	dbConf, err := config.NewDBConf()
	if err != nil {
		return nil, err
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Name)

	if db, err := sqlx.Connect("postgres", dbInfo); err != nil {
		e := erro.NewDBConnError("ConnectToBD", err)
		return nil, e
	} else {
		return db, nil
	}
}

func CreateUser(db *sqlx.DB, tguserID int) erro.Err {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, state) VALUES ($1, $2)", tguserID, "borned")
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("CreateUser", err)
		return e
	}
	return nil
}

func CheckUser(db *sqlx.DB, tguserID int) (bool, erro.Err) {
	var isExist bool
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		e := erro.NewWrapError("CheckUser", err)
		return false, e
	}
	return isExist, nil
}

func ChangeUserState(db *sqlx.DB, tguserID int, state string) erro.Err {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET state=$1 where id=$2", state, tguserID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("CheckUserState", err)
		return e
	}
	return nil
}

func ChangeSelectCategory(db *sqlx.DB, tguserID int, categoryID int) erro.Err {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_category=$1 where id=$2", categoryID, tguserID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("ChangeSelectCategory", err)
		return e
	}
	return nil
}

func ChangeSelectTask(db *sqlx.DB, tguserID int, taskID int) erro.Err {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_task=$1 where id=$2", taskID, tguserID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("ChangeSelectTask", err)
		return e
	}
	return nil
}

func CheckUserState(db *sqlx.DB, tguserID int) (string, erro.Err) {
	var state string
	err := db.QueryRow("SELECT state FROM tguser WHERE id=$1", tguserID).Scan(&state)
	if err != nil {
		e := erro.NewWrapError("CheckUserState", err)
		return "", e
	}
	return state, nil
}

func CheckSelectTaskID(db *sqlx.DB, tguserID int) (int, erro.Err) {
	var taskID int
	err := db.QueryRow("SELECT select_task FROM tguser WHERE id=$1", tguserID).Scan(&taskID)
	if err != nil {
		e := erro.NewWrapError("CheckSelectTaskID", err)
		return 0, e
	}
	return taskID, nil
}

func CheckSelectCategoryID(db *sqlx.DB, tguserID int) (int, erro.Err) {
	var categoryID int
	err := db.QueryRow("SELECT select_category FROM tguser WHERE id=$1", tguserID).Scan(&categoryID)
	if err != nil {
		e := erro.NewWrapError("CheckSelectCategory", err)
		return 0, e
	}
	return categoryID, nil
}

func CreateCategory(db *sqlx.DB, tguserID int, name string) erro.Err {
	var categoryID int
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO category (name) VALUES ($1)", name)
	err := tx.QueryRow("SELECT id FROM category where name=$1", name).Scan(&categoryID)
	if err != nil {
		e := erro.NewWrapError("CreateCategory", err)
		return e
	}
	tx.MustExec("INSERT INTO category_tguser (category_id, tguser_id) VALUES ($1, $2)", categoryID, tguserID)
	tx.Commit()
	return nil
}

func DeleteCategory(db *sqlx.DB, tguserID int, categoryID int) erro.Err {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_category=NULL WHERE id=$1", tguserID)
	tx.MustExec("UPDATE tguser SET select_task=NULL WHERE id=$1", tguserID)
	tx.MustExec("DELETE FROM task WHERE category_id=$1", categoryID)
	tx.MustExec("DELETE FROM category_tguser WHERE category_id=$1", categoryID)
	tx.MustExec("DELETE FROM category WHERE id=$1", categoryID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("DeleteCategory", err)
		return e
	}

	return nil
}

func ListAllCategories(db *sqlx.DB, tguserID int) ([][]string, erro.Err) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM category_tguser WHERE tguser_id=$1)", tguserID).Scan(&isExist)
	if err != nil {
		e := erro.NewWrapError("ListAllCategories", err)
		return nil, e
	}

	if isExist == false {
		return nil, nil
	}

	rows, err := db.Query("SELECT category.id, category.name FROM category LEFT JOIN category_tguser ON category_tguser.category_id=category.id WHERE category_tguser.tguser_id=$1", tguserID)
	if err != nil {
		e := erro.NewWrapError("ListAllCategories", err)
		return nil, e
	}

	allCategories := [][]string{}

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			e := erro.NewWrapError("ListAllCategories", err)
			return nil, e
		}

		allCategories = append(allCategories, []string{strconv.Itoa(id), name})
		defer rows.Close()
	}
	return allCategories, nil
}

func CreateTask(db *sqlx.DB, categoryID int, text string) erro.Err {

	sliceTaskText := strings.SplitN(text, "\n\n", 3)
	var fmtDeadline string
	//Проверяем сколько разных строк пришло
	len := len(sliceTaskText)
	if len == 3 {
		title := sliceTaskText[0]
		deadline := sliceTaskText[1]
		description := sliceTaskText[2]

		//Проверяем не попал ли deadline в title
		titleErr, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, title)
		if titleErr == true {
			err := errors.New("TitleIsNotDeadline")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}
		//Проверяем введен ли title
		if title == "" {
			err := errors.New("NilTitle")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}

		//Проверяем правильность введения deadline
		deadlineOK, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, deadline)
		if deadlineOK == false {
			if deadline == "-" || deadline == "" {
				fmtDeadline = "01.01.1998"

				layout := "02.01.2006"
				t, err := time.Parse(layout, fmtDeadline)
				if err != nil {
					err := errors.New("DateErr")
					e := erro.NewWrapError("CreateTask", err)
					return e
				}
				fmtDeadline = t.Format("01-02-2006")
			}
			err := errors.New("DateErr")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}

		layout := "02.01.2006"
		t, err := time.Parse(layout, deadline)
		if err != nil {
			err := errors.New("DateErr")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}
		fmtDeadline := t.Format("01-02-2006")

		//Проверяем введено ли описание
		if description == "" {
			description = "-"
		}

		tx := db.MustBegin()
		tx.MustExec("INSERT INTO task(title, complete, category_id, description, deadline) VALUES ($1, $2, $3, $4, $5)", title, false, categoryID, description, fmtDeadline)
		tx.Commit()

	} else if len == 2 {
		var deadline string
		var description string

		title := sliceTaskText[0]
		something := sliceTaskText[1]

		//Проверяем не попал ли deadline в title
		titleErr, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, title)
		if titleErr == true {
			err := errors.New("TitleIsNotDeadline")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}
		//Проверяем введен ли title
		if title == "" {
			err := errors.New("NilTitle")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}

		isItDeadline, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, something)
		if isItDeadline == true {
			deadline = something
			description = "-"

			//форматируем deadline
			layout := "02.01.2006"
			t, err := time.Parse(layout, deadline)
			if err != nil {
				err := errors.New("DateErr")
				e := erro.NewWrapError("CreateTask", err)
				return e
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
					err := errors.New("DateErr")
					e := erro.NewWrapError("CreateTask", err)
					return e
				}
				fmtDeadline = t.Format("01-02-2006")

				description = "-"
			}

			description = something
			deadline = "01.01.1998"

			layout := "02.01.2006"
			t, err := time.Parse(layout, deadline)
			if err != nil {
				err := errors.New("DateErr")
				e := erro.NewWrapError("CreateTask", err)
				return e
			}
			fmtDeadline = t.Format("01-02-2006")

		}
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO task(title, complete, category_id, description, deadline) VALUES ($1, $2, $3, $4, $5)", title, false, categoryID, description, fmtDeadline)
		tx.Commit()

	} else if len == 1 {
		title := sliceTaskText[0]
		var deadline string
		var description string
		//Проверяем не попал ли deadline в title
		titleErr, _ := regexp.MatchString(`^\d{2}(\.)\d{2}(\.)\d{4}$`, title)
		if titleErr == true {
			err := errors.New("TitleIsNotDeadline")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}
		//Проверяем введен ли title
		if title == "" {
			err := errors.New("NilTitle")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}

		deadline = "01.01.1998"

		layout := "02.01.2006"
		t, err := time.Parse(layout, deadline)
		if err != nil {
			err := errors.New("DateErr")
			e := erro.NewWrapError("CreateTask", err)
			return e
		}
		fmtDeadline = t.Format("01-02-2006")

		description = "-"
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO task(title, complete, category_id, description, deadline) VALUES ($1, $2, $3, $4, $5)", title, false, categoryID, description, fmtDeadline)
		tx.Commit()
	}
	return nil
}

func ViewTask(db *sqlx.DB, categoryID int, taskID int, tguserID int) (string, erro.Err) {

	var text string

	var title string
	var complete bool
	var description string
	var deadline string

	tx := db.MustBegin()
	err := tx.QueryRow("SELECT title, description, complete, deadline-now()::date as deadline FROM task WHERE category_id=$1 AND id=$2", categoryID, taskID).Scan(&title, &description, &complete, &deadline)
	if err != nil {
		e := erro.NewWrapError("ViewTask", err)
		return "", e
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

	return text, nil
}

func ChangeTaskTitle(db *sqlx.DB, tguserID int, taskID int, text string) erro.Err {
	tx := db.MustBegin()
	db.MustExec("UPDATE task SET title=$1 WHERE id=$2", text, taskID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("ChangeTaskTitle", err)
		return e
	}
	return nil
}

func ChangeTaskDescription(db *sqlx.DB, tguserID int, taskID int, text string) erro.Err {
	tx := db.MustBegin()
	db.MustExec("UPDATE task SET description=$1 WHERE id=$2", text, taskID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("ChangeTaskDescription", err)
		return e
	}
	return nil
}

func ChangeTaskDeadline(db *sqlx.DB, tguserID int, taskID int, text string) erro.Err {
	tx := db.MustBegin()
	db.MustExec("UPDATE task SET deadline=$1 WHERE id=$2", text, taskID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("ChangeTaskDeadline", err)
		return e
	}
	return nil
}

func ListTasks(db *sqlx.DB, categoryID int) (string, erro.Err) {
	var isExist string
	err := db.QueryRow("SELECT exists (SELECT 1 FROM task WHERE category_id=$1)", categoryID).Scan(&isExist)
	if err != nil {
		e := erro.NewWrapError("ListTasks", err)
		return "", e
	}
	tasks := [][]string{}
	if isExist == "true" {
		rows, err := db.Query("SELECT id, title, complete FROM task WHERE category_id=$1 ORDER BY complete", categoryID)
		if err != nil {
			e := erro.NewWrapError("ListTasks", err)
			return "", e
		}

		for rows.Next() {
			var id int
			var title string
			var complete bool
			rows.Scan(&id, &title, &complete)
			tasks = append(tasks, []string{strconv.Itoa(id), title, strconv.FormatBool(complete)})
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

func DeleteTask(db *sqlx.DB, taskID int, tguserID int) erro.Err {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET select_task=NULL WHERE id=$1", tguserID)
	tx.MustExec("DELETE FROM task WHERE id=$1", taskID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("DeleteTask", err)
		return e
	}
	return nil
}

func CompleteTask(db *sqlx.DB, taskID int) erro.Err {
	tx := db.MustBegin()
	tx.MustExec("UPDATE task SET complete=$1 WHERE id=$2", true, taskID)
	err := tx.Commit()
	if err != nil {
		e := erro.NewWrapError("CompleteTask", err)
		return e
	}
	return nil
}

func IsComplete(db *sqlx.DB, taskID int) (bool, erro.Err) {
	var isComplete bool
	tx := db.MustBegin()
	err := tx.QueryRow("SELECT complete FROM task where id=$1", taskID).Scan(&isComplete)
	if err != nil {
		e := erro.NewWrapError("IsComplete", err)
		return false, e
	}
	tx.Commit()
	return isComplete, nil
}

func IsTaskExist(db *sqlx.DB, categoryID int, taskID int) (bool, erro.Err) {
	var isExist bool
	err := db.QueryRow("SELECT exists (SELECT 1 FROM task WHERE id=$1 AND category_id=$2)", taskID, categoryID).Scan(&isExist)
	if err != nil {
		e := erro.NewWrapError("IsTaskExist", err)
		return false, e
	}
	return isExist, nil
}
