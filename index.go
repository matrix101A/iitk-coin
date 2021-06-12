package main

import (
	"log"
	"net/http"

	"github.com/matrix101A/handlers"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

func main() {

	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/secretpage", handlers.SecretPageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
