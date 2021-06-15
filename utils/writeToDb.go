package utils

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func WriteUserToDb(name string, rollno string, password string) error {
	database, _ :=
		sql.Open("sqlite3", "./database/user.db")

	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT,rollno TEXT PRIMARY KEY,password TEXT)")

	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO user (name,rollno,password) VALUES (?, ?, ?)")
	_, err := statement.Exec(name, rollno, password)
	if err != nil {
		return err
	}
	err = InitializeCoins(rollno)
	if err != nil {
		return err
	}
	database.Close()
	return nil

}

func InitializeCoins(rollno string) error {
	database, _ :=
		sql.Open("sqlite3", "./database/user.db")

	statement, _ :=
		database.Prepare(`INSERT INTO bank (rollno,coins) VALUES ($1,$2); `)

	_, err := statement.Exec(rollno, 0)
	if err != nil {
		return err
	}
	database.Close()
	return nil

}

func WriteCoinsToDb(rollno string, numberOfCoins string) error {

	coins_int, e := strconv.Atoi(numberOfCoins)
	if e != nil {
		return e
	}

	database, _ :=
		sql.Open("sqlite3", "./database/user.db")
	_, err := GetUserFromRollNo(rollno)

	if err != nil {
		return err
	}
	total_coins, err := GetCoinsFromRollNo(rollno)
	if err != nil {
		return err
	}
	total_coins = total_coins + coins_int
	statement, _ :=
		database.Prepare(`UPDATE bank SET coins = $1 WHERE rollno= $2;`)
	_, err = statement.Exec(total_coins, rollno)
	if err != nil {
		return err
	}
	database.Close()
	return nil
}

func TransferCoinDb(firstRollno string, secondRollno string, transferAmount int) error {
	if firstRollno == secondRollno {
		return nil
	}
	db, _ := sql.Open("sqlite3", "./database/user.db")
	tx, err := db.Begin()
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return err
	}
	firstUserCoins, err := GetCoinsFromRollNo(firstRollno)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	firstUserCoins = firstUserCoins - transferAmount // withdraw from first user
	if firstUserCoins < 0 {
		balanceError := errors.New("not enough balance ")
		return balanceError
	}
	res, execErr := tx.Exec("UPDATE bank SET coins = ? WHERE rollno= ?;", firstUserCoins, firstRollno)
	rowsAffected, _ := res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		return err
	}
	secondUserCoins, err := GetCoinsFromRollNo(secondRollno)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	secondUserCoins = secondUserCoins + transferAmount // Deposit to second user

	res, execErr = tx.Exec("UPDATE bank SET coins = ? WHERE rollno= ?;", secondUserCoins, secondRollno)
	rowsAffected, _ = res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	db.Close()
	return nil
}
