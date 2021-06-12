package handlers

import (
	"database/sql"
	"encoding/json"
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

func Get_hashed_password(rollno string) string {
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
func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		hashedPassword := Get_hashed_password(rollno)

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
