package utils

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func WriteUserToDb(name string, rollno string, password string) error {
	statement, _ :=
		Db.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT,rollno TEXT PRIMARY KEY,password TEXT)")

	statement.Exec()

	statement, _ =
		Db.Prepare("INSERT INTO user (name,rollno,password) VALUES (?, ?, ?)")
	_, err := statement.Exec(name, rollno, password)
	if err != nil {
		return err
	}
	err = InitializeCoins(rollno)
	if err != nil {
		return err
	}

	return nil

}

func InitializeCoins(rollno string) error {
	Db, _ :=
		sql.Open("sqlite3", "./database/user.db")

	statement, _ :=
		Db.Prepare(`INSERT INTO bank (rollno,coins) VALUES ($1,$2); `)

	_, err := statement.Exec(rollno, 0)
	if err != nil {
		return err
	}

	return nil

}

func WriteCoinsToDb(rollno string, numberOfCoins string) error {

	coins_int, e := strconv.Atoi(numberOfCoins)
	if e != nil {
		return e
	}
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
		Db.Prepare(`UPDATE bank SET coins = $1 WHERE rollno= $2;`)
	_, err = statement.Exec(total_coins, rollno)
	if err != nil {
		return err
	}

	return nil
}

func TransferCoinDb(firstRollno string, secondRollno string, transferAmount int) error {
	if firstRollno == secondRollno {
		return nil
	}
	_, err := GetUserFromRollNo(firstRollno)
	if err != nil {
		return errors.New("user " + firstRollno + " not present ")
	}
	_, err = GetUserFromRollNo(secondRollno)
	if err != nil {
		return errors.New("user " + secondRollno + " not present ")
	}

	var options = sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	tx, err := Db.BeginTx(context.Background(), &options)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return err
	}

	res, execErr := tx.Exec("UPDATE bank SET coins = coins - ? WHERE rollno=? AND coins - ? >= 0", transferAmount, firstRollno, transferAmount)

	rowsAffected, _ := res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return err
		}

		balanceError := errors.New("not enough balance ")
		return balanceError

	}

	res, execErr = tx.Exec("UPDATE bank SET coins = coins + ? WHERE rollno=? ", transferAmount, secondRollno)

	rowsAffected, _ = res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
