package bd

import (
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToBD() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=qrlzvrn dbname=clozapinum sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateUser(db *sqlx.DB, userID int) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, state) VALUES ($1, $2)", userID, "borned")
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func CheckUser(db *sqlx.DB, userID int) error {
	var isExist string
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1)", userID).Scan(&isExist)
	if err != nil {
		return err
	}
	return nil
}

func ChangeUserState(db *sqlx.DB, userID int, state string) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET state=$1 where id=$2", state, userID)
	tx.Commit()
	return nil
}

func CheckUserState(db *sqlx.DB, userID int) (string, error) {
	var state string
	err := db.QueryRow("SELECT state FROM tguser WHERE id=$1", userID).Scan(&state)
	if err != nil {
		return "", err
	}
	return state, nil
}

func CheckSelectTaskID(db *sqlx.DB, userID int) (int, error) {
	var taskID int
	err := db.QueryRow("SELECT select_task FROM tguser WHERE id=$1", userID).Scan(&taskID)
	if err != nil {
		return 0, err
	}
	return taskID, nil
}

func CheckSelectCategoryID(db *sqlx.DB, userID int) (int, error) {
	var categoryID int
	err := db.QueryRow("SELECT select_category FROM tguser WHERE id=$1", userID).Scan(&categoryID)
	if err != nil {
		return 0, err
	}
	return categoryID, nil
}

func CreateCategory(db *sqlx.DB, userID int, name string) error {
	var categoryID int
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO category (name) VALUES ($1)", name)
	err := tx.QueryRow("SELECT id FROM category where name=$1", name).Scan(&categoryID)
	if err != nil {
		return err
	}
	tx.MustExec("INSERT INTO category_tguser (category_id, tguser_id) VALUES ($1, $2)", categoryID, userID)
	tx.Commit()
	return nil
}

func ListAllCategories(db *sqlx.DB, userID int) ([][]string, error) {
	rows, err := db.Query("SELECT category.id, category.name FROM category LEFT JOIN category_tguser ON category_tguser.category_id=category.id WHERE category_tguser.tguser_id=$1", userID)
	if err != nil {
		return nil, err
	}

	allCategories := [][]string{}

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		allCategories = append(allCategories, []string{strconv.Itoa(id), name})
		defer rows.Close()
	}
	return allCategories, nil
}

func CreateTask(db *sqlx.DB, categoryID int, title string) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO task(title, complete, category_id) VALUES ($1, $2, $3)", title, false, categoryID)
	tx.Commit()
	return nil
}

func ChangeTask() error {
	//
	//
	return nil
}

func ListTasks(db *sqlx.DB, category_id int) ([][]string, error) {
	rows, err := db.Query("SELECT id, title, complete FROM task WHERE category_id=$1", category_id)
	if err != nil {
		return nil, err
	}

	tasks := [][]string{}

	for rows.Next() {
		var id int
		var title string
		var complete bool
		rows.Scan(&id, &title, &complete)
		tasks = append(tasks, []string{strconv.Itoa(id), title, strconv.FormatBool(complete)})
		defer rows.Close()
	}
	return tasks, nil
}

func DeleteTask(db *sqlx.DB, taskID int, userID int) error {
	var user int

	tx := db.MustBegin()
	err := tx.QueryRow("SELECT category_tguser.tguser_id FROM task LEFT JOIN category_tguser ON category_tguser.category_id=task.category_id WHERE task.id=$1 AND category_tguser.tguser_id=$2", taskID, userID).Scan(&user)
	if err != nil {
		return err
	}
	tx.MustExec("DELETE FROM task WHERE id=$1", taskID)
	tx.Commit()

	return nil
}

func CompleteTask(db *sqlx.DB, taskID int, userID int) error {
	var user int

	tx := db.MustBegin()
	err := tx.QueryRow("SELECT category_tguser.tguser_id FROM task LEFT JOIN category_tguser ON category_tguser.category_id=task.category_id WHERE task.id=$1 AND category_tguser.tguser_id=$2", taskID, userID).Scan(&user)
	if err != nil {
		return err
	}
	tx.MustExec("UPDATE task set complete=$1 WHERE id=$2", true, taskID)
	tx.Commit()

	return nil
}

func IsComplete(db *sqlx.DB, taskID int, userID int) (bool, error) {
	var user int
	var isComplete bool
	tx := db.MustBegin()
	err := tx.QueryRow("SELECT category_tguser.tguser_id FROM task LEFT JOIN category_tguser ON category_tguser.category_id=task.category_id WHERE task.id=$1 AND category_tguser.tguser_id=$2", taskID, userID).Scan(&user)
	if err != nil {
		return false, err
	}
	err = tx.QueryRow("SELECT complete FROM task where id=$1", taskID).Scan(&isComplete)
	if err != nil {
		return false, err
	}
	tx.Commit()
	return isComplete, nil
}
