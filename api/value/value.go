package value

import (
	"encoding/json"
	"fmt"
	"iot-project/mqtt"
	"iot-project/tools"
	"iot-project/tools/auth"
	valuePackage "iot-project/tools/value"
	"net/http"
)

func init() {
	http.Handle("/api/setValue", tools.ModifyHTTPCors(http.HandlerFunc(setValue)))
	http.Handle("/api/getValues", tools.ModifyHTTPCors(http.HandlerFunc(getValues)))
	http.Handle("/api/getLastValues", tools.ModifyHTTPCors(http.HandlerFunc(getLastValues)))
}

func setValue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "error")
		return
	}

	type setValueForm struct {
		VarId     uint64 `json:"varId"`
		Value     string `json:"value"`
		MqttToken string `json:"mqttToken"`
		Label     string `json:"label"`
	}

	formData := setValueForm{}

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		fmt.Println(err)
		fmt.Println(r.Body)
		fmt.Fprintf(w, "error")
		return
	}
	valuePackage.SaveValue(mqtt.Client, formData.MqttToken, formData.VarId, formData.Value, true, false)
	fmt.Fprintf(w, "{\"errno\": 0}")
}

type value struct {
	Current_Value string
	CreationDate  string
}

type valuelast struct {
	Id    int
	Value string
}

func getValues(w http.ResponseWriter, r *http.Request) {
	varid := r.FormValue("varid")
	rows, _ := tools.DBS.Query("SELECT value, \"createdAt\" FROM values WHERE varid=$1 ORDER BY id asc", varid)
	var values []value
	for rows.Next() {
		var current_value string
		var createdAt string
		rows.Scan(&current_value, &createdAt)
		res1D := value{
			Current_Value: current_value,
			CreationDate:  createdAt}
		values = append(values, res1D)
	}
	if len(values) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		res1B, _ := json.Marshal(values)
		fmt.Fprintf(w, string(res1B))
	}
}

func getLastValues(w http.ResponseWriter, r *http.Request) {
	userId, _ := auth.GetActorIdByRequest(r)
	rows, _ := tools.DBS.Query("select distinct uv.varid, (select v.id from values v where v.varid=uv.varid order by \"createdAt\" desc limit 1), (select v.value from values v where v.varid=uv.varid order by \"createdAt\" desc limit 1) from users_variables uv where uv.userid=$1", userId)
	var values []valuelast
	for rows.Next() {
		var id int
		var val string
		var varid int
		rows.Scan(&varid, &id, &val)
		res1D := valuelast{
			Id:    varid,
			Value: val}
		values = append(values, res1D)
	}
	res1B, _ := json.Marshal(values)
	fmt.Fprintf(w, string(res1B))
}
