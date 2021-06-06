package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func get_hashed_password(rollno string) string {
	database, _ :=
		sql.Open("sqlite3", "./database/user.db")
	rollno_int, _ := strconv.Atoi(rollno)
	sqlStatement := `SELECT password FROM user WHERE rollno= $1;`
	row := database.QueryRow(sqlStatement, rollno_int)

	var hashed_password string
	row.Scan(&hashed_password)
	fmt.Println(hashed_password)
	return (hashed_password)

}
func write_details(name string, rollno string, password string, w http.ResponseWriter) error {
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
	return nil

}
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		name := r.FormValue("name")
		rollno := r.FormValue("rollno")
		password := r.FormValue("password")

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {

			log.Fatal(err)
		}

		write_err := write_details(name, rollno, string(hashed_password), w)
		if write_err != nil {
			log.Printf("Body read error, %v", write_err)

			w.WriteHeader(500) // Return 500 Internal Server Error.
			w.Write([]byte("Roll number must be unique "))

			return
		}
		fmt.Println("Your account was created sucessfully ")
		fmt.Fprintf(w, "Accont created sucessfully")
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func CreateToken(userRollNo string) (string, error, time.Time) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "thisisasecret") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_roll_no"] = userRollNo
	expTime := time.Now().Add(time.Minute * 15)
	atClaims["exp"] = expTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("thisisasecret"))
	if err != nil {
		return "", err, time.Now()
	}
	return token, nil, expTime
}

func VerifyToken(request_token string) (*jwt.Token, error) {
	tokenString := request_token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("thisisasecret"), nil //enter secret key here
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenMetadata(user_token string) (string, error) { //returns the roll no of the user
	token, err := VerifyToken(user_token)
	if err != nil {
		return " ", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		roll_no, _ := claims["user_roll_no"].(string)
		return roll_no, err
	}

	return " ", err

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {

	case "POST":

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			w.WriteHeader(500) // log error
			return
		}
		rollno := r.FormValue("rollno")
		password := r.FormValue("password")
		hashedPassword := get_hashed_password(rollno)
		fmt.Println(hashedPassword)
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			w.WriteHeader(500) // send server error
			w.Write([]byte("wrong Password"))
			return
		}
		token, err, expirationTime := CreateToken(rollno)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}
		http.SetCookie(w, &http.Cookie{ // setting cookie for the user with expiration time
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("Password was correct!, You are logged in "))
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Sorry, only  POST methods are supported.")
	}

}

func secretPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		fmt.Fprint(w, "Error 404 Page not found")
		return

	}
	switch r.Method {
	case "GET":
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		user_roll_no, _ := ExtractTokenMetadata(tokenFromUser)
		fmt.Println(user_roll_no, "Hello")

		fmt.Fprint(w, "Welcome to the secret page "+user_roll_no)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Sorry, only GET  methods are supported.")
	}

}
func main() {
	fmt.Println("yoy")
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/secretpage", secretPageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
