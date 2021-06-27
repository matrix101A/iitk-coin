package utils

import (
	"database/sql"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB
var err error
var MaxCoins float64

func ConnectToDb() error {
	godotenv.Load()
	MaxCoins, _ = strconv.ParseFloat(os.Getenv("MAXCOINS"), 32)
	Db, err =
		sql.Open("sqlite3", "./database/user.db")
	if err != nil {
		return err
	}
	return nil
}
