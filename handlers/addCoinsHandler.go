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
	Rollno string `json:"rollno"`
	Coins  string `json:"coins"`
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
	w.Header().Set("Content-Type", "application/json")

	resp := &serverResponse{
		Message: "",
	}

	switch r.Method {

	case "POST":

		var coinsData Bank

		err := json.NewDecoder(r.Body).Decode(&coinsData)
		if err != nil {
			//fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rollno := coinsData.Rollno

		numberOfCoins := coinsData.Coins

		if rollno == "" {
			w.WriteHeader(401)
			resp.Message = "Please enter a roll number"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		_, err = strconv.Atoi(numberOfCoins)
		if err != nil {
			w.WriteHeader(401)
			resp.Message = "Coins should be valid integer"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		err = utils.WriteCoinsToDb(rollno, numberOfCoins)
		if err != nil {
			w.WriteHeader(500)
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, " -User not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = coinsData.Coins + " Coins added to user " + coinsData.Rollno
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
