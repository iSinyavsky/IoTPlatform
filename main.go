package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"io/ioutil"
	_ "iot-project/api/codeGen"
	_ "iot-project/api/getActor"
	_ "iot-project/api/login"
	_ "iot-project/api/logout"
	_ "iot-project/api/metrics"
	_ "iot-project/api/register"
	_ "iot-project/api/test"
	_ "iot-project/api/value"
	_ "iot-project/api/variable"
	database "iot-project/db"
	_ "iot-project/http-api"
	"iot-project/mqtt"
	"iot-project/tools"
	_ "iot-project/tools/triggers"
	_ "iot-project/tools/yandex"
	_ "iot-project/websocket/connect"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {

	var err error

	rand.Seed(time.Now().UnixNano())
	config := make(map[string]string)

	content, err := ioutil.ReadFile("iotProject" + ".conf")
	_ = pq.Efatal
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	if err := json.Unmarshal(content, &config); err != nil {
		panic(err)
	}
	connect := config["DBConnectionString"]
	if connect == "" {
		log.Println("No DBConnectionString record in config file")
		return
	}

	tools.DBS, err = sql.Open("postgres", connect)
	if err != nil {
		panic(err)
	}
	defer tools.DBS.Close()

	x, err := strconv.ParseInt(config["MaxOpenConnections"], 10, 32)
	if err != nil {
		x = 256
	}
	tools.DBS.SetMaxOpenConns(int(x))

	x, err = strconv.ParseInt(config["MaxIdleConnections"], 10, 32)
	if err != nil {
		x = 64
	}
	tools.DBS.SetMaxIdleConns(int(x))

	x, err = strconv.ParseInt(config["MaxSessionsFromClient"], 10, 32)
	if err != nil {
		x = 7
	}

	tools.MaxSessionsFromClient = int(x)

	x, err = strconv.ParseInt(config["ServerID"], 10, 32)
	if err != nil || x < 0 {
		log.Println("No suitable ServerID record in config file")
		return
	}
	tools.ServerID = int32(x)

	x, err = strconv.ParseInt(config["MaxUploadFileSize"], 10, 64)
	if err != nil || x < 0 {
		log.Println("No suitable MaxUploadFileSize record in config file")
		return
	}
	tools.MaxUploadFileSize = x

	x, err = strconv.ParseInt(config["MaxUploadedFilesSize"], 10, 64)
	if err != nil || x < 0 {
		log.Println("No suitable MaxUploadedFilesSize record in config file")
		return
	}
	tools.MaxUploadedFilesSize = x

	x, err = strconv.ParseInt(config["MaxUploadedFiles"], 10, 64)
	if err != nil || x < 0 {
		log.Println("No suitable MaxUploadedFiles record in config file")
		return
	}
	tools.MaxUploadedFiles = x

	prod, err := strconv.ParseInt(config["IsProduction"], 10, 32)
	if err != nil || prod < 0 || prod > 1 {
		log.Println("No suitable IsProduction record (value 0 or 1) in config file")
		return
	}
	tools.IsProduction = prod != 0
	tools.ServerAddr = config["ServerAddr"]
	tools.ServerPort = config["ServerPort"]

	tools.DBS.SetConnMaxLifetime(time.Hour)

	database.MigrateDB(tools.DBS, "")
	mqtt.StartMQTT()
	log.Println("Server started")

	log.Fatal(http.ListenAndServe(":1234", nil))

}
