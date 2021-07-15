// person can pick items from a list costing of diiferent coins to redeems it, wubsequent coins will be deducted from the user
// Currently I have a predefined table with three items with corresponding ids of. The redeemed item will be added to users table to reflect the same
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matrix101A/utils"
	_ "github.com/mattn/go-sqlite3"
)

type respondRedeem struct {
	RequestId int    `json:"requestid"`
	Action    string `json:"action"`
}

func RespondRedeemHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/respondredeem" {
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
			http.Error(w, "User not logged in ", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	_, Acctype, _ := utils.ExtractTokenMetadata(tokenFromUser)

	if Acctype != "admin" {
		http.Error(w, "Unauthorized!! admins are allowed ", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp := &serverResponse{
		Message: "",
	}
	switch r.Method {

	case "POST":
		var requestData respondRedeem
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requestId := requestData.RequestId
		action := requestData.Action
		msg, err := utils.RespondRedeemDb(requestId, action)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp.Message = err.Error()
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		resp.Message = msg
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
