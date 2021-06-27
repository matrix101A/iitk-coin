package utils

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func CreateToken(userRollNo string, userType string) (string, time.Time, error) {
	var err error
	//Creating Access Token

	godotenv.Load()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["accountType"] = userType
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
