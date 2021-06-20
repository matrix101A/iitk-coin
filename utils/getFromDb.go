package utils

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func Get_hashed_password(rollno string) string {

	rollno_int, _ := strconv.Atoi(rollno)
	sqlStatement := `SELECT password FROM user WHERE rollno= $1;`
	row := Db.QueryRow(sqlStatement, rollno_int)

	var hashed_password string
	row.Scan(&hashed_password)

	return (hashed_password)

}

func GetCoinsFromRollNo(rollno string) (int, error) {

	statement, _ :=
		Db.Prepare("CREATE TABLE IF NOT EXISTS bank (rollno TEXT PRIMARY KEY ,coins INT)")
	statement.Exec()

	sqlStatement := `SELECT coins FROM bank WHERE rollno= $1;`
	row := Db.QueryRow(sqlStatement, rollno)

	var coins int
	err := row.Scan(&coins)

	if err != nil {
		return 0, err
	}
	return coins, nil

}

func GetUserFromRollNo(rollno string) (*sql.Row, error) {

	sqlStatement := `SELECT name FROM user WHERE rollno= $1;`
	row := Db.QueryRow(sqlStatement, rollno)
	err := row.Scan(&rollno)

	if err != nil {
		return nil, err
	}
	return row, nil
}
