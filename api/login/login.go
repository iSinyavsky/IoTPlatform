package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"iot-project/tools/auth"
	"net/http"
	"strings"
	"time"

	"iot-project/tools"
)

func init() {
	http.Handle("/api/login", tools.ModifyHTTPCors(http.HandlerFunc(handleFormLogin)))
}

type formLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func validateFormLoginData(formData *formLoginData) error {
	email := strings.TrimSpace(formData.Email)
	password := strings.TrimSpace(formData.Password)
	if email == "" || password == "" {
		return errors.New("empty params")
	}
	//if !tools.ValidateEmail(email) || len([]rune(password)) < 6 {
	//	return errors.New("Not validate")
	//}
	return nil
}

func handleFormLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "BAD method")
		return
	}
	formData := formLoginData{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	err = validateFormLoginData(&formData)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	email := strings.TrimSpace(formData.Email)
	password := strings.TrimSpace(formData.Password)
	session, err := loginForm(email, password)
	fmt.Println(session)
	if err != nil {
		fmt.Fprintf(w, "{\"errno\": 1}")
		return
	}
	sessionCookieExpirationTime := time.Now().Add(auth.SESSION_TOKEN_EXPIRATION)
	auth.SetSessionToken(w, session.Token, sessionCookieExpirationTime)
	fmt.Fprintf(w, "{\"errno\": 0, \"userId\": %d}", session.ActorId)
}

func loginForm(email, password string) (*auth.UserSession, error) {

	var actorId uint32
	fmt.Println(auth.HashPassword(password))
	err := tools.DBS.QueryRow("SELECT id FROM users WHERE email = $1 AND password = $2", email, auth.SHA256DigestString(password)).Scan(&actorId)
	token := auth.CreateAuthToken(actorId)
	if token == "" {
		err := fmt.Errorf("Bad token")
		return nil, err
	}
	session := &auth.UserSession{
		Token:   token,
		ActorId: actorId,
	}
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Not found")
	}
	return session, nil
}
