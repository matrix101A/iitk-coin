// person can pick items from a list costing of diiferent coins to redeems it, wubsequent coins will be deducted from the user
// Currently I have a predefined table with three items with corresponding ids of. The redeemed item will be added to users table to reflect the same
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/matrix101A/utils"
	_ "github.com/mattn/go-sqlite3"
)

type redeemCoinsData struct {
	Item_id int `json:"itemid"`
}

func RedeemCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/redeem" {
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
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	rollno, _, _ := utils.ExtractTokenMetadata(tokenFromUser)
	w.Header().Set("Content-Type", "application/json")

	resp := &serverResponse{
		Message: "",
	}

	switch r.Method {

	case "POST":

		var redeemData redeemCoinsData

		err := json.NewDecoder(r.Body).Decode(&redeemData)
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		item_id := redeemData.Item_id

		if rollno == "" {
			w.WriteHeader(401)
			resp.Message = "Please enter a roll number"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		coins, err := utils.RedeemCoinsDb(rollno, item_id) // withdraw from first user and transfer to second
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		resp.Message = "Your resquest is awaiting confirmatino for item " + fmt.Sprintf("%d", item_id) + " .Coins remaining after redeem will be  " + fmt.Sprintf("%.2f", coins)
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
