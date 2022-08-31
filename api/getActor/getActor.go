package getActor

import (
	"database/sql"
	"fmt"
	"iot-project/tools"
	"iot-project/tools/auth"
	"net/http"
)

func init() {
	http.Handle("/api/getActor", tools.ModifyHTTPCors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := auth.GetSessionToken(r)
		if err != nil {
			fmt.Println("aaa1", err)
			fmt.Fprintf(w, "{\"errno\": 1}")
			return

		}

		userId, err := auth.GetActorIdByAuthToken(token)
		if err != nil {
			fmt.Println("aaa2", err)
			fmt.Fprintf(w, "{\"errno\": 1}")
			return
		}

		var mqttToken string
		var nullName sql.NullString
		var name string
		var email string
		var yaTokenNull sql.NullString
		var yaToken string

		tools.DBS.QueryRow("SELECT token, name, email, ya_token FROM users WHERE id = $1", userId).Scan(&mqttToken, &nullName, &email, &yaTokenNull)

		name = nullName.String
		if !nullName.Valid {
			name = "IoT user"
		}
		if yaTokenNull.Valid {
			yaToken = yaTokenNull.String
		}
		fmt.Fprintf(w, "{\"errno\": 0, \"userId\": %d, \"mqttToken\": \"%s\", \"email\": \"%s\", \"name\": \"%s\", \"yaToken\": \"%s\"}", userId, mqttToken, email, name, yaToken)
	})))
}
