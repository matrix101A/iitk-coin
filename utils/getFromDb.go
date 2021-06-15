package utils

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func Get_hashed_password(rollno string) string {
	database, _ :=
		sql.Open("sqlite3", "./database/user.db")
	rollno_int, _ := strconv.Atoi(rollno)
	sqlStatement := `SELECT password FROM user WHERE rollno= $1;`
	row := database.QueryRow(sqlStatement, rollno_int)

	var hashed_password string
	row.Scan(&hashed_password)
	//fmt.Println(hashed_password)
	return (hashed_password)

}

func GetCoinsFromRollNo(rollno string) (int, error) {

	database, _ :=
		sql.Open("sqlite3", "./database/user.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS bank (rollno TEXT PRIMARY KEY ,coins INT)")
	statement.Exec()

	sqlStatement := `SELECT coins FROM bank WHERE rollno= $1;`
	row := database.QueryRow(sqlStatement, rollno)

	var coins int
	err := row.Scan(&coins)
	//fmt.Println(hashed_password)
	if err != nil {
		return 0, err
	}
	return coins, nil

}

func GetUserFromRollNo(rollno string) (*sql.Row, error) {
	database, _ :=
		sql.Open("sqlite3", "./database/user.db")

	sqlStatement := `SELECT name FROM user WHERE rollno= $1;`
	row := database.QueryRow(sqlStatement, rollno)
	err := row.Scan(&rollno)
	//fmt.Println(hashed_password)
	if err != nil {
		return nil, err
	}
	return row, nil
}
