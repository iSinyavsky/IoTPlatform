package register

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"iot-project/tools"
	"iot-project/tools/auth"
	//	"database/sql"
	//	_ "github.com/lib/pq"
)

func init() {
	http.Handle("/api/register", tools.ModifyHTTPCors(http.HandlerFunc(confirm)))
}

type registerFormData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func confirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// fmt.Fprintf(w, tools.EResult(tools.ECodeBadHTTPMethod))
		fmt.Fprintf(w, "Bad request Method")
		return
	}
	formData := registerFormData{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		fmt.Println(err)
		return
	}
	email := formData.Email
	password := formData.Password

	session, err := writeConfirmedUserToDB(email, password)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	sessionCookieExpirationTime := time.Now().Add(auth.SESSION_TOKEN_EXPIRATION)
	auth.SetSessionToken(w, session.Token, sessionCookieExpirationTime)
	fmt.Fprintf(w, "{\"errno\": 0, \"userId\": %d}", session.ActorId)
}

func writeConfirmedUserToDB(email string, password string) (*auth.UserSession, error) {
	userExist := 0
	tools.DBS.QueryRow("SELECT count(id) FROM users WHERE email = $1", email).Scan(&userExist)

	var userId uint32
	if userExist == 0 {
		tools.DBS.QueryRow("INSERT INTO users(email, password, token) VALUES($1, $2, $3) RETURNING id", email, auth.SHA256DigestString(password), auth.SHA256DigestString(email)).Scan(&userId)
	} else {
		return nil, errors.New("User exists")
	}
	return auth.CreateSession(userId)
}
