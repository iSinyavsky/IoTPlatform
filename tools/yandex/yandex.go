package yandex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"iot-project/mqtt"
	"iot-project/tools"
	"iot-project/tools/auth"
	"iot-project/tools/value"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var actorToYandexHTTPClient sync.Map
var actorTokens sync.Map

type State struct {
	Value interface{} `json:"value"`
}
type Capability struct {
	State State  `json:"state"`
	Type  string `json:"type"`
}
type Device struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	Capabilities []Capability `json:"capabilities"`
}
type DevicesResponse struct {
	Status  string   `json:"status"`
	Devices []Device `json:"devices"`
}

func init() {
	http.Handle("/api/yandex/saveYandexToken", tools.ModifyHTTPCors(http.HandlerFunc(saveYandexToken)))
	http.Handle("/api/yandex/getDevices", tools.ModifyHTTPCors(http.HandlerFunc(getDevices)))
	//http.Handle("/api/yandex/testAction", tools.ModifyHTTPCors(http.HandlerFunc(testAction)))
	http.Handle("/api/yandex/sendValue", tools.ModifyHTTPCors(http.HandlerFunc(sendValue)))
	http.Handle("/api/yandex/getDevicesActions", tools.ModifyHTTPCors(http.HandlerFunc(getDevicesActions)))
	http.Handle("/api/yandex/getDevice", tools.ModifyHTTPCors(http.HandlerFunc(getDevice)))
}
func GetYandexHTTPClient(actorId uint32) (*http.Client, string) {
	fmt.Print("auth")
	jar, err := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	payload := strings.NewReader("login=sinyavskiyivan1%40yandex.ru&passwd=3008MAKn")

	req, err := http.NewRequest("POST", "https://passport.yandex.ru/passport?mode=auth&retpath=https://yandex.ru", payload)

	if err != nil {
		fmt.Println(err)
		return nil, ""
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	req, err = http.NewRequest("GET", "https://yandex.ru/quasar/iot", payload)
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, ""
	}

	r, _ := regexp.Compile(`"csrfToken2":"(.+?)"`)
	csrf := r.FindStringSubmatch(string(body))[1]
	actorToYandexHTTPClient.Store(1, client)
	actorTokens.Store(1, csrf)
	return client, csrf
}

func getAuth(actorId int) (*http.Client, string) {
	clientInterface, ok := actorToYandexHTTPClient.Load(actorId)
	tokenInterface, ok := actorTokens.Load(actorId)
	if !ok {
		clientInterface, tokenInterface = GetYandexHTTPClient(uint32(actorId))
	}

	client := clientInterface.(*http.Client)
	token := tokenInterface.(string)
	return client, token
}

func getDevicesActions(w http.ResponseWriter, r *http.Request) {
	actions := value.GetYandexActions(false)
	json, _ := json.Marshal(actions)
	fmt.Fprintf(w, string(json))
}
func saveYandexToken(w http.ResponseWriter, r *http.Request) {
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

	type saveYandexData struct {
		Token interface{} `json:"token"`
	}
	formData := saveYandexData{}
	json.NewDecoder(r.Body).Decode(&formData)

	fmt.Println("userId", userId, formData.Token)

	tools.DBS.Exec("UPDATE users SET ya_token = $1 WHERE id = $2", formData.Token, userId)
	fmt.Fprintf(w, string(""))
}

//func sendValue(w http.ResponseWriter, r *http.Request) {
//
//	actions := GetActions(false)
//	actionType := "devices.capabilities.on_off"
//	currentAction := actions["devices.types.light"][actionType]
//	currentAction.Type = actionType
//	deviceId := "5865b60a-620b-443d-96f4-8e46ab50e292"
//	json, _ := json.Marshal(currentAction)
//
//	fmt.Fprintln(w, string(json))
//	payload := strings.NewReader(`{"actions":[`+string(json)+`]}`)
//	req, err := http.NewRequest("POST", "https://iot.quasar.yandex.ru/m/user/devices/"+deviceId+"/actions", payload)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	req.Header.Add("Content-Type", "text/plain")
//
//	client, token := getAuth(1)
//
//	req.Header.Add("x-csrf-token", token)
//
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Fprint(w, string(body))
//
//}

func sendValue(w http.ResponseWriter, r *http.Request) {

	actorId, _ := auth.GetActorIdByRequest(r)
	deviceId, _ := tools.GetStringParamFromRequestUrl(r, "deviceID")
	varID, _ := tools.GetIntParamFromRequestUrl(r, "varID")
	deviceId = strings.Split(deviceId, "-n")[0]
	instanceType, _ := tools.GetStringParamFromRequestUrl(r, "type")
	deviceType, _ := tools.GetStringParamFromRequestUrl(r, "deviceType")
	val, _ := tools.GetStringParamFromRequestUrl(r, "value")
	var finalValue interface{}
	finalValue = val
	if instanceType == "devices.capabilities.on_off" {
		finalValue, _ = strconv.ParseBool(val)
	}
	if instanceType == "devices.capabilities.range" {
		finalValue, _ = strconv.ParseInt(val, 10, 32)
	}
	actionsList := value.GetYandexActions(finalValue)
	action := actionsList[deviceType][instanceType]
	out, _ := json.Marshal(action)
	payload := strings.NewReader(`{"actions":[` + string(out) + `]}`)
	req, err := http.NewRequest("POST", "https://api.iot.yandex.net/v1.0/devices/"+deviceId+"/actions", payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "text/plain")

	client := http.Client{}

	req.Header.Add("Authorization", "Bearer AQAAAAAd_dTZAAfKlDiEDBf58kMRpqkTbFPzf00")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	mqttToken := ""
	tools.DBS.QueryRow("SELECT token FROM users WHERE id = $1", actorId).Scan(&mqttToken)

	value.SaveValue(mqtt.Client, mqttToken, uint64(varID), stateValueToPlatformValue(finalValue), true, false)

	fmt.Fprint(w, string(body))
}

func stateValueToPlatformValue(value interface{}) string {
	if value == nil {
		return ""
	}
	if reflect.TypeOf(value).Kind() == reflect.String {
		return value.(string)
	}
	if reflect.TypeOf(value).Kind() == reflect.Bool {
		boolToIntString := "0"
		if value.(bool) == true {
			boolToIntString = "1"
		}
		return boolToIntString
	}
	if reflect.TypeOf(value).Kind() == reflect.Float64 {
		return fmt.Sprintf("%f", value.(float64))
	}
	return ""
}
func getDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	req, err := http.NewRequest("GET", "https://api.iot.yandex.net/v1.0/user/info", strings.NewReader(""))
	req.Header.Add("Authorization", "Bearer AQAAAAAd_dTZAAfKlDiEDBf58kMRpqkTbFPzf00")
	if err != nil {
		fmt.Println(err)
		return
	}

	actorId, _ := auth.GetActorIdByRequest(r)
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	deviceResponse := DevicesResponse{}
	err = json.Unmarshal(body, &deviceResponse)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	capTypes := map[string]string{
		"devices.capabilities.on_off":        "On_Off",
		"devices.capabilities.range":         "Range",
		"devices.capabilities.mode":          "Mode",
		"devices.capabilities.color_setting": "Color",
	}
	mqttToken := ""
	tools.DBS.QueryRow("SELECT token FROM users WHERE id = $1", actorId).Scan(&mqttToken)
	tools.DBS.Query("DELETE FROM variables USING users_variables WHERE integration_service = 'yandex' AND users_variables.userid = $1", actorId)
	for _, device := range deviceResponse.Devices {
		for k, capability := range device.Capabilities {
			capType := capTypes[capability.Type]
			name := device.Name + ": " + capType
			val := stateValueToPlatformValue(capability.State.Value)
			varId := 0
			err = tools.DBS.QueryRow("INSERT INTO variables(name, label, integration_service, integration_id, integration_type, style, integration_capability) VALUES($1, $2, 'yandex', $2, $3, $4, $5) RETURNING id", name, fmt.Sprintf("%s-n%d", device.Id, k), device.Type, "{\"bg\": \"#131F39\", \"icon\": \"fa-link\"}", capability.Type).Scan(&varId)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}

			_, err = tools.DBS.Exec("INSERT INTO users_variables(userid, varid) VALUES($1, $2)", actorId, varId)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			value.SaveValue(mqtt.Client, mqttToken, uint64(varId), val, true, false)
		}

	}

	fmt.Fprint(w, string(body))
}

func getDevice(w http.ResponseWriter, r *http.Request) {
	deviceID, _ := tools.GetStringParamFromRequestUrl(r, "deviceID")
	deviceID = strings.Split(deviceID, "-n")[0]
	req, err := http.NewRequest("GET", "https://api.iot.yandex.net/v1.0/devices/"+deviceID, strings.NewReader(""))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer AQAAAAAd_dTZAAfKlDiEDBf58kMRpqkTbFPzf00")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	type State struct {
		Instance string      `json:"instance"`
		Value    interface{} `json:"value"`
	}
	type Capability struct {
		State State `json:"state"`
	}
	type DeviceResponse struct {
		Capabilities []Capability `json:"capabilities"`
	}

	fmt.Fprint(w, string(body))
}
