package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/matrix101A/utils"
	_ "github.com/mattn/go-sqlite3"
)

type Bank struct {
	Rollno  string `json:"rollno"`
	Coins   string `json:"coins"`
	Remarks string `json:"remarks"`
}

func AddCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/addcoins" {
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	_, Acctype, _ := utils.ExtractTokenMetadata(tokenFromUser)

	if Acctype == "member" {
		http.Error(w, "Unauthorized!! Only CTM and admins are allowed ", http.StatusUnauthorized)
		return
	}

	resp := &serverResponse{
		Message: "",
	}

	switch r.Method {

	case "POST":

		var coinsData Bank

		err := json.NewDecoder(r.Body).Decode(&coinsData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rollno := coinsData.Rollno

		numberOfCoins := coinsData.Coins

		remarks := coinsData.Remarks

		if rollno == "" {
			w.WriteHeader(401)
			resp.Message = "Please enter a roll number"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		_, userAccType, _ := utils.GetUserFromRollNo(rollno)
		if userAccType == "CTM" && Acctype == "CTM" {
			http.Error(w, "Unauthorized only admins are alowed ", http.StatusUnauthorized)
			return
		}
		if userAccType == "admin" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		_, err = strconv.ParseFloat(numberOfCoins, 32)
		if err != nil {
			w.WriteHeader(401)
			resp.Message = "Coins should be valid number "
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		err, errorMessage := utils.WriteCoinsToDb(rollno, numberOfCoins, remarks)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, errorMessage)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = errorMessage + coinsData.Coins + " Coins added to user " + coinsData.Rollno
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
