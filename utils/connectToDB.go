package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB
var err error

func ConnectToDb() error {
	Db, err =
		sql.Open("sqlite3", "./database/user.db")
	if err != nil {
		return err
	}
	return nil
}
