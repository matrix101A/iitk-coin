package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ :=
		sql.Open("sqlite3", "./database/student.db")

	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT,rollno INTEGER PRIMARY KEY)")

	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO user (name,rollno) VALUES (?, ?)")
	statement.Exec("Abhinav Tiwari", 290031)
	rows, _ :=
		database.Query("SELECT name,rollno FROM user")
	var rollno int
	var name string

	for rows.Next() {
		rows.Scan(&name, &rollno)
		fmt.Println(strconv.Itoa(rollno) + ": " + name)
	}
}
