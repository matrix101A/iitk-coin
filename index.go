package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/matrix101A/handlers"
	"github.com/matrix101A/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/secretpage", handlers.SecretPageHandler)
	http.HandleFunc("/addcoins", handlers.AddCoinsHandler)
	http.HandleFunc("/transfercoin", handlers.TransferCoinHandler)
	http.HandleFunc("/getcoins", handlers.GetCoinsHandler)
	http.HandleFunc("/redeem", handlers.RedeemCoinsHandler)
	http.HandleFunc("/additems", handlers.AddItemsHandler)
	http.HandleFunc("/pendingredeems", handlers.PendingRedeemsHandler)
	http.HandleFunc("/respondredeem", handlers.RespondRedeemHandler)

	err := utils.ConnectToDb()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Database connection sucessful")
	log.Fatal(http.ListenAndServe(":8080", nil))
	defer utils.Db.Close()

}
