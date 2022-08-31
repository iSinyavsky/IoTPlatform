package http_api

import (
	"encoding/json"
	"fmt"
	"iot-project/mqtt"
	"iot-project/tools"
	"iot-project/tools/value"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Value struct {
	Value string `json:"value"`
	Time  string `json:"time"`
}

func init() {
	http.Handle("/http-api/", http.HandlerFunc(handleFormLogin))
}

func handleFormLogin(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	urlParts := strings.Split(url, "/")
	token := urlParts[2]
	if len(urlParts) < 4 {
		fmt.Fprintf(w, "{\"errno\": \"Не правильная ссылка\"}")
		return
	}
	label := urlParts[3]
	action := ""
	if len(urlParts) == 5 {
		action = urlParts[4]
	}
	fmt.Println(action)
	userId := 0
	tools.DBS.QueryRow("SELECT id FROM users u WHERE token = $1", token).Scan(&userId)
	if userId == 0 {
		fmt.Fprintf(w, "{\"errno\": \"Отказ в доступе\"}")
		return
	}

	varId := 0
	tools.DBS.QueryRow("SELECT va.id FROM users u INNER JOIN users_variables uv ON uv.userid = u.id INNER JOIN variables va ON va.label = $1 WHERE token = $2", label, token).Scan(&varId)

	if varId == 0 {
		fmt.Fprintf(w, "{\"errno\": \"Отказ в доступе\"}")
		return
	}

	if r.Method == "GET" {
		id := 0
		values := make([]Value, 0)
		filters, errString := applyFilters(r)
		if errString != "" {
			fmt.Fprintf(w, "{\"errno\": \"%s\"}", errString)
			return
		}
		rows, _ := tools.DBS.Query("SELECT * FROM values WHERE varid = $1 "+filters, varId)
		count := 0
		for rows.Next() {
			value := Value{}
			varId := 0
			rows.Scan(&id, &varId, &value.Value, &value.Time)
			values = append(values, value)
		}

		valuesBytes, _ := json.Marshal(values)
		fmt.Fprintf(w, "{\"errno\": 0, \"count\": %d, \"data\": %s}", count, valuesBytes)

		//fmt.Fprintf(w, string(valuesBytes))
	}

	type valueData struct {
		Value interface{} `json:"value"`
	}
	if r.Method == "POST" {
		formData := valueData{}
		err := json.NewDecoder(r.Body).Decode(&formData)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		val := ""
		if reflect.TypeOf(formData.Value).Kind() == reflect.Float64 {
			val = fmt.Sprintf("%f", formData.Value.(float64))
		} else if reflect.TypeOf(formData.Value).Kind() == reflect.String {
			val = formData.Value.(string)
		}

		value.SaveValue(mqtt.Client, token, uint64(varId), val, true, false)

		fmt.Fprintf(w, "{\"errno\": 0}")
	}
}

func applyFilters(r *http.Request) (string, string) {
	sql := ""
	lt, _ := tools.GetStringParamFromRequestUrl(r, "lt")
	le, _ := tools.GetStringParamFromRequestUrl(r, "le")
	eq, _ := tools.GetStringParamFromRequestUrl(r, "eq")
	ne, _ := tools.GetStringParamFromRequestUrl(r, "ne")
	ge, _ := tools.GetStringParamFromRequestUrl(r, "ge")
	gt, _ := tools.GetStringParamFromRequestUrl(r, "gt")
	limit, _ := tools.GetStringParamFromRequestUrl(r, "limit")
	order, _ := tools.GetStringParamFromRequestUrl(r, "order")

	if lt != "" {
		err := validateDate(lt)
		if err != nil {
			return "", "Неправильный формат даты"
		}
		sql += " AND \"createdAt\" < '" + lt + "'"
	}
	if le != "" {
		err := validateDate(le)
		if err != nil {
			return "", "Неправильный формат даты"
		}
		sql += " AND \"createdAt\" <= '" + le + "'"
	}
	if eq != "" {
		err := validateDate(eq)
		if err != nil {
			return "", "Неправильный формат даты"
		}
		sql += " AND \"createdAt\" = '" + eq + "'"
	}
	if ne != "" {
		err := validateDate(ne)
		if err != nil {
			fmt.Println(err)
			return "", "Неправильный формат даты"
		}
		sql += " AND \"createdAt\" != '" + ne + "'"
	}
	if ge != "" {
		err := validateDate(ge)
		if err != nil {
			return "", "Неправильный формат даты"
		}
		sql += " AND \"createdAt\" >= '" + ge + "'"
	}
	if gt != "" {
		err := validateDate(gt)
		if err != nil {
			return "", "Неправильный формат даты"
		}
		sql += " AND \"createdAt\" > '" + gt + "'"
	}
	if order != "" && (order == "desc" || order == "asc") {
		sql += " ORDER BY \"createdAt\" " + order
	}
	if limit != "" {
		sql += " LIMIT " + limit
	}

	return sql, ""
}

func validateDate(date string) error {
	_, err := time.Parse(time.RFC3339, date)
	if err != nil {
		_, err = time.Parse("2006-01-02", date)
		if err != nil {
			_, err = time.Parse("2006-01-02T15:04", date)
			if err != nil {
				_, err = time.Parse("2006-01-02 15:04", date)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return nil
}
