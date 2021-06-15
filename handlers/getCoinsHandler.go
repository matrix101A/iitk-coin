package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/matrix101A/utils"
	_ "github.com/mattn/go-sqlite3"
)

func GetCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/getcoins" {
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

	case "GET":

		var user User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			//fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Rollno == "" {
			w.WriteHeader(401)
			resp.Message = "Please enter a roll number"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		coins, err := utils.GetCoinsFromRollNo(user.Rollno)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Fprintf(w, " -User not found")
			return
		}

		w.WriteHeader(http.StatusOK)
		resp.Message = "Your coins are " + strconv.Itoa(coins)
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Sorry, only GET requests are supported"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
