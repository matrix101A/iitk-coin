package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func WriteToDb(name string, rollno string, password string) error {
	database, _ :=
		sql.Open("sqlite3", "./database/user.db")

	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT,rollno TEXT PRIMARY KEY,password TEXT)")

	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO user (name,rollno,password) VALUES (?, ?, ?)")
	_, err := statement.Exec(name, rollno, password)
	if err != nil {
		return err
	}
	return nil

}
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	resp := &serverResponse{
		Message: "",
	}
	switch r.Method {

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		/* if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		} */
		var user User

		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			//fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//fmt.Println(user)

		name := user.Name
		rollno := user.Rollno
		password := user.Password
		if rollno == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)

			resp.Message = "Roll No or Password Cannot be empty"

			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			//log.Fatal(err)
			w.WriteHeader(401)

			resp.Message = "Server error"

			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
		}

		write_err := WriteToDb(name, rollno, string(hashed_password))

		if write_err != nil {
			//log.Printf("Body read error, %v", write_err)

			w.WriteHeader(500) // Return 500 Internal Server Error.
			resp.Message = "Roll no already exists"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}

		//fmt.Println("Your account was created sucessfully ")

		w.WriteHeader(http.StatusOK)
		//Write json response back to response

		resp.Message = "Your account has benn created. To login, go to /login"

		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)

		resp.Message = "Only POST requests are supported"

		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
}
