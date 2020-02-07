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

func CreateUser(db *sqlx.DB, id int) error {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO tguser (id, state) VALUES ($1, $2)", id, "borned")
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func CheckUser(db *sqlx.DB, id int) error {
	var isExist string
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1);", id).Scan(&isExist)
	if err != nil {
		return err
	}
	return nil
}

func ChangeUserState(db *sqlx.DB, id int, state string) error {
	tx := db.MustBegin()
	tx.MustExec("UPDATE tguser SET state=$1 where id=$2", state, id)
	tx.Commit()
	return nil
}

func CreateCategory(db *sqlx.DB, userID int, name string) error {
	var categoryID int
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO category (name) VALUES ($1)", name)
	err := tx.QueryRow("SELECT id FROM category where name = $1", name).Scan(&categoryID)
	if err != nil {
		return err
	}
	tx.MustExec("INSERT INTO category_tguser (category_id, tguser_id) VALUES ($1, $2)", categoryID, userID)
	tx.Commit()
	return nil
}

func ListAllCategory(db *sqlx.DB, userID int) error {
	rows, err := db.Query("SELECT category.id, category.name FROM category LEFT JOIN category_tguser ON category_tguser.category_id=category.id WHERE category_tguser.tguser_id=$1", userID)
	if err != nil {
		return err
	}

	allCategory := [][]string{}

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}

		allCategory = append(allCategory, []string{strconv.Itoa(id), name})
		defer rows.Close()
	}
	return nil
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
