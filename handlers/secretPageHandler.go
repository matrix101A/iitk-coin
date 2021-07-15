package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matrix101A/utils"
	_ "github.com/mattn/go-sqlite3"
)

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
				resp.Message = "User not logged in"
				JsonRes, _ := json.Marshal(resp)
				w.Write(JsonRes)

				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		user_roll_no, Acctype, err := utils.ExtractTokenMetadata(tokenFromUser)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			resp.Message = "Access Unauthorized "
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		resp.Message = "Welcome to the secret page " + user_roll_no + " " + Acctype
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
