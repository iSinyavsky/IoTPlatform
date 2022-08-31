package value

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"iot-project/tools"
	"iot-project/websocket/notifier"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type IfEvent struct {
	Value    string `json:"value"`
	Operator string `json:"operator"`
	Type     int    `json:"type"`
}

type ThenEvent struct {
	Value string `json:"value"`
}

func comparing(value1 string, value2 string, compare string) {

}

type StateType struct {
	Instance string       `json:"instance"`
	Value    *interface{} `json:"value,omitempty"`
}

type Action struct {
	Type  string    `json:"type"`
	State StateType `json:"state"`
}

func GetYandexActions(value interface{}) map[string]map[string]Action {
	actions := make(map[string]map[string]Action, 0)

	actions["devices.types.light"] = make(map[string]Action, 0)
	actions["devices.types.light"]["devices.capabilities.on_off"] = Action{Type: "devices.capabilities.on_off", State: StateType{Instance: "on", Value: &value}}
	actions["devices.types.light"]["devices.capabilities.color_setting"] = Action{Type: "devices.capabilities.color_setting", State: StateType{Value: &value}}
	actions["devices.types.light"]["devices.capabilities.range"] = Action{Type: "devices.capabilities.range", State: StateType{Instance: "brightness", Value: &value}}
	actions["devices.types.light"]["devices.capabilities.toggle"] = Action{Type: "devices.capabilities.toggle", State: StateType{Instance: "backlight"}}

	return actions
}

func SendValueToYandex(actorId uint32, varID int, val string) {

	instanceType := ""
	deviceType := ""
	deviceId := ""
	integrationService := ""

	tools.DBS.QueryRow("SELECT integration_service, integration_capability, integration_type, integration_id FROM variables WHERE id = $1", varID).Scan(&integrationService, &instanceType, &deviceType, &deviceId)

	if integrationService == "" {
		return
	}
	deviceId = strings.Split(deviceId, "-n")[0]

	fmt.Println("value???", varID, deviceId)

	var finalValue interface{}
	finalValue = val
	if instanceType == "devices.capabilities.on_off" {
		finalValue, _ = strconv.ParseBool(val)
	}
	if instanceType == "devices.capabilities.range" {
		finalValue, _ = strconv.ParseInt(val, 10, 32)
	}
	actionsList := GetYandexActions(finalValue)
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

	fmt.Println(string(body))

}

func SaveValue(client mqtt.Client, token string, variableId uint64, value string, withPublish bool, withInsert bool) {
	actorId := 0
	tools.DBS.QueryRow("SELECT userid FROM users_variables WHERE varid = $1 LIMIT 1", variableId).Scan(&actorId)
	fmt.Println("save")
	if withInsert {
		_, err := tools.DBS.Exec("INSERT INTO values (value, varid) VALUES($1, $2)", value, variableId)

		if err != nil {
			fmt.Println(err)
		}
	}
	now := time.Now()

	eventsRows, _ := tools.DBS.Query("SELECT if_event, then_event FROM events WHERE if_event -> $1 is not null", variableId)

	for eventsRows.Next() {
		ifEventString := ""
		thenEventString := ""

		fmt.Println("event")
		eventsRows.Scan(&ifEventString, &thenEventString)
		ifEvents := make(map[uint64]IfEvent, 0)
		thenEvents := make(map[uint64]ThenEvent, 0)

		ifTrue := true

		err := json.Unmarshal([]byte(ifEventString), &ifEvents)
		if err != nil {
			fmt.Println(err)
			fmt.Println("fuck")
		}
		err = json.Unmarshal([]byte(thenEventString), &thenEvents)

		if err != nil {
			fmt.Println(err)
		}

		for varID, ifEvent := range ifEvents {
			compareValue := ""
			tools.DBS.QueryRow("SELECT value FROM values WHERE varid = $1 ORDER BY id DESC LIMIT 1", varID).Scan(&compareValue)
			if !compare(compareValue, ifEvent.Value, ifEvent.Operator) {
				ifTrue = false
			}
		}

		if ifTrue {
			for variableId, thenEvent := range thenEvents {
				var label string
				tools.DBS.QueryRow("SELECT label FROM variables WHERE id = $1", variableId).Scan(&label)
				//tx.Exec("INSERT INTO values (value, varid) VALUES($1, $2)", thenEvent.Value, variableId)
				//tx.Rollback()
				client.Publish(fmt.Sprintf("value/%s/%s", token, label), 0, false, thenEvent.Value)
				SendValueToYandex(uint32(actorId), int(variableId), thenEvent.Value)
			}
		}

	}
	if withPublish {
		label := ""
		err := tools.DBS.QueryRow("SELECT label FROM variables WHERE id = $1", variableId).Scan(&label)
		if err != nil {
			fmt.Println("errrrroooor", err)
			return
		}
		client.Publish(fmt.Sprintf("value/%s/%s", token, label), 0, false, value)
	}
	fmt.Println("tut?")

	notifier.NotifyByActorId(uint32(actorId), fmt.Sprintf("{\"type\":\"newValue\", \"varId\": %d, \"date\": %d, \"value\": \"%s\"}", variableId, now.UnixNano()/1000, value))
}

func compare(a string, b string, operator string) bool {
	fmt.Println("compare", a, operator, b)
	if operator == "&lt;" {
		if a < b {
			return true
		}
	}

	if operator == "&le;" {
		if a <= b {
			return true
		}
	}

	if operator == "=" {
		if a == b {
			return true
		}
	}

	if operator == "&ge;" {
		if a >= b {
			return true
		}
	}
	if operator == "&gt;" {
		if a > b {
			return true
		}
	}

	return false
}
