package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matrix101A/utils"
	_ "github.com/mattn/go-sqlite3"
)

type transferCoin struct {
	Account_1_roll_no string `json:"firstrollno"`
	Account_2_roll_no string `json:"secondrollno"`
	Amount            int    `json:"amount"`
}

func TransferCoinHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/transfercoin" {
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

		var transferData transferCoin

		err := json.NewDecoder(r.Body).Decode(&transferData)
		if err != nil {
			//fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		firstRollno := transferData.Account_1_roll_no
		secondRollno := transferData.Account_2_roll_no
		transferAmount := transferData.Amount

		if firstRollno == "" || secondRollno == "0" {
			w.WriteHeader(401)
			resp.Message = "Please enter a roll number"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		err = utils.TransferCoinDb(firstRollno, secondRollno, transferAmount) // withdraw from first user and transfer to second
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = "Transaction Sucessfull !"
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
