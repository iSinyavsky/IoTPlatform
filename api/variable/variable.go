package variable

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"iot-project/tools"
	"iot-project/tools/auth"
	"iot-project/tools/value"
	"net/http"
	"time"
)

func init() {
	http.Handle("/api/getVariables", tools.ModifyHTTPCors(http.HandlerFunc(getVariables)))
	http.Handle("/api/addVariable", tools.ModifyHTTPCors(http.HandlerFunc(addVariable)))
	http.Handle("/api/deleteVariable", tools.ModifyHTTPCors(http.HandlerFunc(deleteVariable)))
	http.Handle("/api/updateVariable", tools.ModifyHTTPCors(http.HandlerFunc(updateVariable)))
	http.Handle("/api/updateStyleVariable", tools.ModifyHTTPCors(http.HandlerFunc(updateStyleVariable)))
	http.Handle("/api/events", tools.ModifyHTTPCors(http.HandlerFunc(addEvent)))
	http.Handle("/api/getEvents", tools.ModifyHTTPCors(http.HandlerFunc(getEvents)))
	http.Handle("/api/removeEvent", tools.ModifyHTTPCors(http.HandlerFunc(removeEvent)))
	http.Handle("/api/removeVariable", tools.ModifyHTTPCors(http.HandlerFunc(removeVariable)))
}

type addForm struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type variable struct {
	Id          int
	Name        string
	Label       string
	ServiceName string
	IntType     string
	IntID       string
	Style       json.RawMessage
	Capability  string
}

type event struct {
	If_event   string
	Then_event string
}

type Trigger struct {
	Time     string
	Interval int
}

type event_response struct {
	Id           int
	If_event     map[int]value.IfEvent
	If_trigger   Trigger
	Then_event   map[int]value.ThenEvent
	RunTime      string
	CreationDate string
	IsActive     bool
}

func getVariables(w http.ResponseWriter, r *http.Request) {
	userId, _ := auth.GetActorIdByRequest(r)
	rows, _ := tools.DBS.Query("SELECT v.id, v.name, v.label, v.integration_service, v.integration_type, v.integration_id, v.style, v.integration_capability FROM variables v, users_variables uv WHERE v.id=uv.varid AND uv.userid=$1", userId)
	var variables []variable
	for rows.Next() {
		var id int
		var name string
		var label string
		var service sql.NullString
		var integrationId sql.NullString
		var integrationType sql.NullString
		var style sql.NullString
		var integrationCapability string
		rows.Scan(&id, &name, &label, &service, &integrationType, &integrationId, &style, &integrationCapability)

		res1D := variable{
			Id:          id,
			Name:        name,
			Label:       label,
			ServiceName: service.String,
			IntID:       integrationId.String,
			IntType:     integrationType.String,
			Capability:  integrationCapability,
		}

		if style.Valid {
			res1D.Style = json.RawMessage(style.String)
		}
		variables = append(variables, res1D)
	}
	res1B, _ := json.Marshal(variables)
	fmt.Fprintf(w, string(res1B))
}

func removeVariable(w http.ResponseWriter, r *http.Request) {
	id, _ := tools.GetIntParamFromRequestUrl(r, "id")
	tools.DBS.Exec("DELETE FROM variables WHERE id = $1", id)
	tools.DBS.Exec("DELETE FROM users_variables WHERE varid = $1", id)
	fmt.Fprintf(w, "{\"ok\": 1}")
}
func addVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Method Not Allowed")
	}
	userId, _ := auth.GetActorIdByRequest(r) //получить id пользователя

	fmt.Println(userId)
	formData := addForm{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		fmt.Println(err)
		fmt.Println(r.Body)
		fmt.Fprintf(w, "error")
		return
	}

	fmt.Println(formData.Name)
	var Varid uint64
	tools.DBS.QueryRow("INSERT INTO variables(name,label) VALUES($1,$2) RETURNING id", formData.Name, formData.Label).Scan(&Varid)
	tools.DBS.Exec("INSERT INTO users_variables(userid,varid) VALUES($1,$2)", userId, Varid)

	fmt.Fprintf(w, "{\"ok\": 1}")
}

func deleteVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "Method Not Allowed")
	}
	id := r.FormValue("id")
	tools.DBS.Exec("DELETE FROM variables WHERE id=$1", id)
	fmt.Fprintf(w, "{\"ok\": 1}")
}

func updateVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		fmt.Fprintf(w, "Method Not Allowed")
	}
	id := r.FormValue("id")
	rows, _ := tools.DBS.Query("SELECT name, label FROM variables WHERE id=$1", id)
	var name string
	var label string
	rows.Scan(&name, &label)
	formData := addForm{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		fmt.Println(err)
		fmt.Println(r.Body)
		fmt.Fprintf(w, "error")
		return
	}
	tools.DBS.Exec("UPDATE variables SET name=$1,label=$2 WHERE id=$3", formData.Name, formData.Label, id)
	fmt.Fprintf(w, "{\"ok\": 1}")
}

func updateStyleVariable(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Method Not Allowed")
	}
	id := r.FormValue("id")

	type StyleStruct struct {
		Icon string `json:"icon"`
		Bg   string `json:"bg"`
	}

	var p StyleStruct

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	st, err := json.Marshal(p)
	fmt.Println("aa", id, string(st), err)

	tx, _ := tools.DBS.Begin()
	_, err = tx.Exec("UPDATE variables SET style=$1 WHERE id=$2", fmt.Sprintf("%s", st), id)
	if err != nil {
		tx.Rollback()
		fmt.Println("uu", err)
		return
	}
	tx.Commit()
	fmt.Fprintf(w, "{\"ok\": 1}")
}

func addEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Method Not Allowed")
	}

	formData := event_response{}

	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "error")
		return
	}
	fmt.Println("heello")

	var runTime sql.NullTime
	var interval sql.NullInt64
	if formData.If_trigger.Time != "" {
		runTime.Time, _ = time.Parse(time.RFC3339, formData.If_trigger.Time)
		runTime.Valid = true
	}
	if formData.If_trigger.Interval != 0 {
		interval.Int64 = int64(formData.If_trigger.Interval)
		interval.Valid = true
		runTime.Time = time.Now().Add(time.Minute * time.Duration(formData.If_trigger.Interval))
		runTime.Valid = true
	}
	fmt.Println(runTime.Time)
	ifEventString, _ := json.Marshal(formData.If_event)
	thenEventString, _ := json.Marshal(formData.Then_event)

	tools.DBS.Exec("INSERT INTO events(if_event,then_event, runtime, interval) VALUES($1,$2,$3,$4)", string(ifEventString), string(thenEventString), runTime, interval)

	fmt.Fprintf(w, "{\"ok\": 1}")
}

func removeEvent(w http.ResponseWriter, r *http.Request) {
	id, _ := tools.GetIntParamFromRequestUrl(r, "id")
	tools.DBS.Exec("DELETE FROM events WHERE id = $1", id)
	fmt.Fprintf(w, "{\"ok\": 1}")
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	rows, err := tools.DBS.Query("SELECT * FROM events ORDER BY (isactive is true), id desc")
	if err != nil {
		fmt.Println(err)
		return
	}
	var events []event_response
	for rows.Next() {
		var id int
		var if_eventString string
		var then_eventString string
		var if_event map[int]value.IfEvent
		var then_event map[int]value.ThenEvent
		var createdAt string
		var isActive bool
		var interval sql.NullInt64
		var runTime sql.NullString
		err := rows.Scan(&id, &if_eventString, &then_eventString, &createdAt, &isActive, &runTime, &interval)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal([]byte(if_eventString), &if_event)
		json.Unmarshal([]byte(then_eventString), &then_event)
		res1D := event_response{
			Id:           id,
			If_event:     if_event,
			Then_event:   then_event,
			If_trigger:   Trigger{Interval: int(interval.Int64), Time: runTime.String},
			CreationDate: createdAt,
			IsActive:     isActive,
			RunTime:      runTime.String,
		}
		events = append(events, res1D)
	}

	if len(events) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		res1B, _ := json.Marshal(events)
		//res := string(regexp.MustCompile(`[\[\\\]]`).ReplaceAll([]byte(res1B), []byte("")))
		fmt.Fprintf(w, string(res1B))
	}
}
