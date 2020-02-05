package bd

import (
	"log"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

func ConnectToBD() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=qrlzvrn dbname=clozapinum sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CreateUser(db *sqlx.DB, id int) error {
	stmt, err := db.Prepare("INSERT INTO tguser(id, state) VALUES ($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(id, "main")
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(lastId)

	defer db.Close()
	return nil
}

func CheckUser(db *sqlx.DB, id int) error {
	var isExist string
	err := db.QueryRow("SELECT exists (select 1 from tguser where id=$1);", id).Scan(&isExist)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func ChangeUserState(db *sqlx.DB, id int) error {
	return nil
}

func CreateCategory() error {

	return nil
}

func ListAllCategory() error {
	//
	//
	return nil
}

func CreateTask() error {
	//
	//
	return nil
}

func ChangeTask() error {
	//
	//
	return nil
}

func DeleteTask() error {
	//
	//
	return nil
}

func CompleteTask() error {
	//
	//
	return nil
}
