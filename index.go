package main

import (
	"log"
	"net/http"

	"github.com/matrix101A/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/secretpage", handlers.SecretPageHandler)
	http.HandleFunc("/addcoins", handlers.AddCoinsHandler)
	http.HandleFunc("/transfercoin", handlers.TransferCoinHandler)
	http.HandleFunc("/getcoins", handlers.GetCoinsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
