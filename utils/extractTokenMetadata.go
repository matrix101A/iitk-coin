package utils

import jwt "github.com/dgrijalva/jwt-go"

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
