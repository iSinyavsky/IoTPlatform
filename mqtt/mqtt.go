package mqtt

import (
	"fmt"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"iot-project/tools"
	"iot-project/tools/value"
)

const server = "mosca-broker:1883"
const username = "mosquitto"
const password = ""

var Client mqtt.Client

func connect(clientId string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", server))
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("connect to mqtt")
	return client
}

var MqttClient mqtt.Client

func listen(topic string) {
	Client = connect("sub")
	Client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		parts := strings.Split(msg.Topic(), "/")
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))

		userId, varId := tools.CheckVariableOrCreate(parts[1], parts[2])

		if userId == 0 || varId == 0 {
			return
		}

		value.SaveValue(client, parts[1], varId, string(msg.Payload()), false, true)

		//events
		//ifvariable, elsevariable, result, type (simple 0, time 0), compare(<,<=,=,>=,=)
	})
}

func StartMQTT() {

	go listen("value/#")

}
