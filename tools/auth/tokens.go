package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"iot-project/tools"
)

type UnconfirmedUserClaims struct {
	*jwt.StandardClaims
	Name  string `json:"uname,omitempty"`
	Email string `json:"email,omitempty"`
	Lang  string `json:"lang,omitempty"`
}

var CONFIRMATION_TOKEN_AUDIENCE = fmt.Sprintf("[%s/api/register/request_confirm, %s/api/register/confirm]", tools.ServerAddr, tools.ServerAddr)

func CreateConfirmationToken(name, email, lang string, expirationTime time.Duration) string {
	expiresAt := time.Now().Add(expirationTime).Unix()
	claims := jwt.MapClaims{
		"uname": name,
		"email": email,
		"lang":  lang,
		"exp":   expiresAt,
		"aud":   CONFIRMATION_TOKEN_AUDIENCE,
	}
	return CreateToken(claims)
}

func ParseConfirmationToken(tokenString string) (*UnconfirmedUserClaims, error) {
	token, err := ParseToken(tokenString, &UnconfirmedUserClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UnconfirmedUserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT token")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("JWT token expired")
	}
	if !claims.VerifyAudience(CONFIRMATION_TOKEN_AUDIENCE, true) {
		return nil, fmt.Errorf("JWT token Audience %s is invalid (must be %s)", claims.Audience, CONFIRMATION_TOKEN_AUDIENCE)
	}
	return claims, nil
}

type RestorePasswordClaims struct {
	*jwt.StandardClaims
	Email string `json:"email,omitempty"`
}

var RESTORE_PASSWORD_TOKEN_AUDIENCE = fmt.Sprintf("[%s/api/password/restore, %s/api/password/commit_restore]", tools.ServerAddr, tools.ServerAddr)

func CreateRestorePasswordToken(email string, expirationTime time.Duration) string {
	expiresAt := time.Now().Add(expirationTime).Unix()
	claims := jwt.MapClaims{
		"email": email,
		"exp":   expiresAt,
		"aud":   RESTORE_PASSWORD_TOKEN_AUDIENCE,
	}
	return CreateToken(claims)
}

func ParseRestorePasswordToken(tokenString string) (*RestorePasswordClaims, error) {
	token, err := ParseToken(tokenString, &RestorePasswordClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*RestorePasswordClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT token")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("JWT token expired")
	}
	if !claims.VerifyAudience(RESTORE_PASSWORD_TOKEN_AUDIENCE, true) {
		return nil, fmt.Errorf("JWT token Audience %s is invalid (must be %s)", claims.Audience, RESTORE_PASSWORD_TOKEN_AUDIENCE)
	}
	return claims, nil
}

type InvitationClaims struct {
	*jwt.StandardClaims
	Lang    string `json:"lang,omitempty"`
	ActorId uint32 `json:"actorId,omitempty"`
}

var INVITE_TOKEN_AUDIENCE = fmt.Sprintf("[%s/api/login/invite]", tools.ServerAddr)

func CreateInvitationToken(lang string, actorid uint32, expirationTime time.Duration) string {
	expiresAt := time.Now().Add(expirationTime).Unix()
	claims := jwt.MapClaims{
		"lang":    lang,
		"actorid": actorid,
		"exp":     expiresAt,
		"aud":     INVITE_TOKEN_AUDIENCE,
	}
	return CreateToken(claims)
}

func ParseInvitationToken(tokenString string) (*InvitationClaims, error) {
	token, err := ParseToken(tokenString, &InvitationClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*InvitationClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT token")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("JWT token expired")
	}
	if !claims.VerifyAudience(INVITE_TOKEN_AUDIENCE, true) {
		return nil, fmt.Errorf("JWT token Audience %s is invalid (must be %s)", claims.Audience, INVITE_TOKEN_AUDIENCE)
	}
	return claims, nil
}

func readGooglePublicRSAKeys() (map[string]*rsa.PublicKey, error) {
	// Google's public RSA keys in PEM format.
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	keysPem := make(map[string]string)
	err = json.Unmarshal(body, &keysPem)
	if err != nil {
		return nil, err
	}

	// Parse keys in PEM format and get *rsa.PublicKey values.
	keys := make(map[string]*rsa.PublicKey)
	for k, v := range keysPem {
		rsaPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(v))
		if err != nil {
			continue
		}
		keys[k] = rsaPubKey
	}
	// TODO Store parsed keys instead of reparsing them every time verifyGoogleJWT is called.
	return keys, nil
}

type GoogleClaims struct {
	*jwt.StandardClaims
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Name          string `json:"name,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
	Lang          string `json:"locale,omitempty"`
	Picture       string `json:"picture,omitempty"`
}

// Workaround for JWTs issued in the future (due to different clock setup), not
// checking iat field shouldn't affect security.
func (c *GoogleClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	return true
}

func ParseGoogleJWT(tokenString string) (*jwt.Token, error) {
	// Retrieve Goggle's public RSA keys.
	rsaPubKeys, err := readGooglePublicRSAKeys()
	if err != nil {
		return nil, err
	}
	// Parse JWT in tokenString. Note that iat check is ignored, see GoogleClaims
	// definition below for more details.
	token, err := jwt.ParseWithClaims(tokenString, &GoogleClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check signing method. For better security there is only one and it is hardcoded.
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Get kid field from JWT header. It is used for picking the correct public key.
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("header field kid is required in order to pick public key")
		}
		// Get the correct public key from `rsaPubKeys map[string]*rsa.PublicKey`.
		rsaPubKey, ok := rsaPubKeys[kid]
		if !ok {
			return nil, fmt.Errorf("public key specified by kid: %s not found", kid)
		}
		return rsaPubKey, nil
	})
	if err != nil {
		return nil, err
	}
	// Check if issuer is valid.
	issuer := token.Claims.(*GoogleClaims).StandardClaims.Issuer
	const ISSUER = "accounts.google.com"
	const ISSUER_HTTPS = "https://accounts.google.com"
	if issuer != ISSUER && issuer != ISSUER_HTTPS {
		return nil, fmt.Errorf("issuer %v is invalid", issuer)
	}
	// Check if auidience is valid.
	audience := token.Claims.(*GoogleClaims).StandardClaims.Audience
	//const GOOGLE_AUDIENCE = "954230887202-9jn2pgrhm5s91n68sle995pgl8at605k.apps.googleusercontent.com"
	const GOOGLE_AUDIENCE = "208549480891-1daccg0ppqr0snhl10o16fobdedpt3af.apps.googleusercontent.com"
	if audience != GOOGLE_AUDIENCE {
		return nil, fmt.Errorf("audience %v is invalid", audience)
	}
	// Check if JWT is valid and return it.
	if _, ok := token.Claims.(*GoogleClaims); ok && token.Valid {
		return token, nil
	} else {
		return nil, err
	}
}

func CreateToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(authTokenBlowFish))
	data := strings.Split(tokenString, ".")
	return data[1] + "." + data[2]
}

func ParseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." + tokenString
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Check signing method. For better security there is only one and it is hardcoded.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authTokenBlowFish), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token, nil
	}
	return nil, err
}

type FacebookClaims struct {
	Lang, Email, Name, Id string
}

type BotUserClaims struct {
	*jwt.StandardClaims
	BotType int    `json:"bottype,omitempty"`
	UserId  string `json:"userid,omitempty"`
}

var BOT_TOKEN_AUDIENCE = fmt.Sprintf("[%s/api/add_bot_user]", tools.ServerAddr)

func CreateBotUserToken(botType int, botUserId string, expirationTime time.Duration) string {
	expiresAt := time.Now().Add(expirationTime).Unix()
	claims := jwt.MapClaims{
		"bottype": botType,
		"userid":  botUserId,
		"exp":     expiresAt,
		"aud":     BOT_TOKEN_AUDIENCE,
	}
	return CreateToken(claims)
}

func ParseBotUserToken(tokenString string) (*BotUserClaims, error) {
	token, err := ParseToken(tokenString, &BotUserClaims{})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*BotUserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid JWT token")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("JWT token expired")
	}
	if !claims.VerifyAudience(BOT_TOKEN_AUDIENCE, true) {
		return nil, fmt.Errorf("JWT token Audience %s is invalid (must be %s)", claims.Audience, CONFIRMATION_TOKEN_AUDIENCE)
	}
	return claims, nil
}
