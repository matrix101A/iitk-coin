package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

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

func SecretPageHandler(w http.ResponseWriter, r *http.Request) {
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
