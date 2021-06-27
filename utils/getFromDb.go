package utils

import (
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

func GetCoinsFromRollNo(rollno string) (float64, error) {

	statement, _ :=
		Db.Prepare("CREATE TABLE IF NOT EXISTS bank (rollno TEXT PRIMARY KEY ,coins INT)")
	statement.Exec()

	sqlStatement := `SELECT coins FROM bank WHERE rollno= $1;`
	row := Db.QueryRow(sqlStatement, rollno)

	var coins float64
	err := row.Scan(&coins)

	if err != nil {
		return 0, err
	}
	return coins, nil

}

func GetUserFromRollNo(rollno string) (string, string, error) {

	sqlStatement := `SELECT name,account_type FROM user WHERE rollno= $1;`
	row := Db.QueryRow(sqlStatement, rollno)
	var userName string
	var userType string
	err := row.Scan(&userName, &userType)
	if err != nil {
		return "", "", err
	}

	return userName, userType, nil
}

func getItemFromId(item_id int) (float64, int, error) {
	var cost float64
	var available int

	sqlStatement := `SELECT cost,available FROM items WHERE id= $1;`
	row := Db.QueryRow(sqlStatement, strconv.Itoa(item_id))

	err := row.Scan(&cost, &available)
	if err != nil {
		return 0, 0, err
	}
	return cost, available, nil
}
