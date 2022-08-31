package codeGen

import (
	"encoding/json"
	"fmt"
	"html/template"
	"iot-project/tools"
	"net/http"
	"strings"
)

type genForm struct {
	Title         string   `json:"title"`
	SSID          string   `json:"ssid"`
	Pass          string   `json:"pass"`
	ModulesChosed []module `json:"modulesChosed"`
	MqttToken     string   `json:"mqttToken"`
}

type module struct {
	Id     string   `json:"id"`
	Label  string   `json:"label"`
	Name   string   `json:"name"`
	Pins   []string `json:"pins"`
	Type   string   `json:"Type"`
	PubSub int      `json:"pubSub"`
	Data   string   `json:"data"`
}

func init() {

	http.Handle("/api/codeGen", tools.ModifyHTTPCors(http.HandlerFunc(codeGenerate)))
}
func codeGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "BAD method")
		return
	}
	formData := genForm{}
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		fmt.Println("ee", err)
		fmt.Println(r.Body)
		return
	}

	tmpl, err := template.ParseFiles("./api/codeGen/index.txt")

	if err != nil {
		fmt.Println(err)
	}
	var b strings.Builder
	fmt.Println(tmpl)
	err = tmpl.Execute(&b, formData)

	if err != nil {
		fmt.Println("aa", err)
	}

	fmt.Println(b.String())

	//fmt.Fprintf(w, b.String())

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+"allEmployees2020.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.Write([]byte(b.String()))
}
