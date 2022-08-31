package login

import (
	"fmt"
	"net/http"

	"iot-project/tools"
	"iot-project/tools/auth"
)

func init() {
	http.Handle("/api/logout", tools.ModifyHTTPCors(http.HandlerFunc(logout)))
}

func logout(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetSessionToken(r)
	if err != nil {
		fmt.Fprintf(w, "{\"errno\": 1, \"message\" : \"User is not logged in\"}")
		return
	}

	auth.UnsetSessionToken(w)

	fmt.Fprintf(w, "{\"errno\": 0}")
	return
}
