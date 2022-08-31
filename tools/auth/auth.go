package auth

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const authTokenBlowFish = "l0$t.fa7f.seCret.fb29.4U"

func CreateAuthToken(actorid uint32) string {
	now := time.Now().Unix()
	arg := fmt.Sprintf("%08x%08x", now, actorid)
	hash := SHA256DigestString(arg + authTokenBlowFish)
	return arg + hash
}

const SESSION_TOKEN_EXPIRATION = 30 * 24 * time.Hour

// Returns actorid and creation time
func GetActorIdByAuthToken(token string) (uint32, error) {
	fmt.Println(len(token), SHA256DigestString(token[0:16]+authTokenBlowFish), token[16:])
	if len(token) != 80 || SHA256DigestString(token[0:16]+authTokenBlowFish) != token[16:] {
		return 0, fmt.Errorf("session token is invalid")
	}
	var unixCreatedAt int64
	unixCreatedAt, err := strconv.ParseInt(token[0:8], 16, 32)
	if err != nil {
		return 0, err
	}
	var actorid uint64
	actorid, err = strconv.ParseUint(token[8:16], 16, 32)
	if err != nil {
		return 0, err
	}
	createdAt := time.Unix(unixCreatedAt, 0)
	expirationTime := createdAt.Add(SESSION_TOKEN_EXPIRATION)
	if time.Now().After(expirationTime) {
		return 0, fmt.Errorf("session token is expired")
	}
	return uint32(actorid), nil
}

func GetActorIdByRequest(r *http.Request) (uint32, error) {
	token, err := GetSessionToken(r)
	if err != nil {
		return 0, fmt.Errorf("session token is invalid")
	}

	actorId, err := GetActorIdByAuthToken(token)
	if err != nil {
		return 0, fmt.Errorf("session token is invalid")
	}

	return actorId, nil
}

const SESSION_COOKIE_NAME = "SessionToken"

// Sets SessionToken cookie with appropriate security options and expirationTime
func SetSessionToken(w http.ResponseWriter, sessionToken string, expirationTime time.Time) {
	cookie := http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    sessionToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expirationTime,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func UnsetSessionToken(w http.ResponseWriter) {
	SetSessionToken(w, "", time.Now())
}

// Get SessionToken from the Cookie request header
func GetSessionToken(r *http.Request) (string, error) {
	if r == nil {
		return "", fmt.Errorf("couldn't get %s cookie: request == nil", SESSION_COOKIE_NAME)
	}
	cookie, err := r.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		fmt.Println("ee", err)
		return "", err
	}
	fmt.Println("norm", cookie.Value)
	return cookie.Value, nil
}

func CreateSession(actorid uint32) (*UserSession, error) {
	token := CreateAuthToken(actorid)
	if token == "" {
		return nil, fmt.Errorf("could not create auth token for user %v", actorid)
	}
	session := &UserSession{token, actorid}
	return session, nil
}

type UserSession struct {
	Token   string
	ActorId uint32
}

func MD5DigestString(s string) string {
	return fmt.Sprintf("%032x", md5.Sum([]byte(s)))
}

func SHA256DigestString(s string) string {
	return fmt.Sprintf("%064x", sha256.Sum256([]byte(s)))
}

func HashPassword(password string) string {
	salt := fmt.Sprintf("%08x", rand.Uint32())
	return fmt.Sprintf("%s%s", salt, SHA256DigestString(salt+password))
}

func CheckPassword(password, hash string) bool {
	if len(hash) != 72 {
		return false
	}
	return hash[8:] == SHA256DigestString(hash[0:8]+password)
}
