package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Rollno   string `json:"rollno"`
	Password string `json:"password"`
}

type serverResponse struct {
	Message string `json:"message"`
}

func get_hashed_password(rollno string) string {
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
func WriteToDb(name string, rollno string, password string) error {
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
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	resp := &serverResponse{
		Message: "",
	}
	switch r.Method {

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		/* if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		} */
		var user User

		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			//fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//fmt.Println(user)

		name := user.Name
		rollno := user.Rollno
		password := user.Password
		if rollno == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "Roll No or Password Cannot be empty"

			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			//log.Fatal(err)
			w.WriteHeader(401)

			resp.Message = "Server error"

			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
		}

		write_err := WriteToDb(name, rollno, string(hashed_password))

		if write_err != nil {
			//log.Printf("Body read error, %v", write_err)

			w.WriteHeader(500) // Return 500 Internal Server Error.
			resp.Message = "Roll no already exists"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		//fmt.Println("Your account was created sucessfully ")

		w.WriteHeader(http.StatusOK)
		//Write json response back to response

		resp.Message = "Your account has benn created. To login, go to /login"

		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Only POST requests are supported"

		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
}

func CreateToken(userRollNo string) (string, time.Time, error) {
	var err error
	//Creating Access Token

	godotenv.Load()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_roll_no"] = userRollNo
	expTime := time.Now().Add(time.Minute * 15)
	atClaims["exp"] = expTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESSKEY")))
	if err != nil {
		return "", time.Now(), err
	}
	return token, expTime, err
}

func VerifyToken(request_token string) (*jwt.Token, error) {
	tokenString := request_token
	godotenv.Load()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESSKEY")), nil //enter secret key here
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
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	resp := &serverResponse{
		Message: "",
	}
	switch r.Method {

	case "POST":

		var user User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			//fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rollno := user.Rollno
		password := user.Password
		hashedPassword := get_hashed_password(rollno)

		if hashedPassword == "" {
			w.WriteHeader(401)
			resp.Message = "User does not exist"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			w.WriteHeader(500) // send server error
			resp.Message = "Password was incorrect"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		token, expirationTime, err := CreateToken(rollno)
		if err != nil {
			w.WriteHeader(401)
			resp.Message = "Server Error"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return

		}

		http.SetCookie(w, &http.Cookie{ // setting cookie for the user with expiration time
			Name:     "token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
		})

		w.WriteHeader(http.StatusOK)

		resp.Message = "Password Correct, you are logged in (Cookie set)"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Sorry, only POST requests are supported"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}

func secretPageHandler(w http.ResponseWriter, r *http.Request) {
	resp := &serverResponse{
		Message: "",
	}
	if r.URL.Path != "/secretpage" {
		w.WriteHeader(404)
		resp.Message = "404 Page not formed"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				resp.Message = "Access restricted, user not authorized"
				JsonRes, _ := json.Marshal(resp)
				w.Write(JsonRes)

				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		user_roll_no, err := ExtractTokenMetadata(tokenFromUser)
		//fmt.Println(user_roll_no, "Hello")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			resp.Message = "Access Unauthorized "
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		resp.Message = "Welcome to the secret page " + user_roll_no
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Only GET requests are supported "
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
	}

}
func main() {

	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/secretpage", secretPageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
