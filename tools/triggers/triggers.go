package value

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"iot-project/mqtt"
	"iot-project/tools"
	"iot-project/tools/value"
	"time"
)

func init() {
	go checkTriggers()
}

func checkTriggers() {
	for {
		time.Sleep(10 * time.Second)

		nowTime := time.Now()
		eventsRows, _ := tools.DBS.Query("SELECT if_event, then_event, id, interval FROM events WHERE runtime < $1 AND isactive = false", nowTime)

		for eventsRows.Next() {
			ifEventString := ""
			thenEventString := ""
			eventId := 0
			var interval sql.NullInt32
			fmt.Println("trigger!!")
			eventsRows.Scan(&ifEventString, &thenEventString, &eventId, &interval)
			thenEvents := make(map[uint64]value.IfEvent, 0)
			err := json.Unmarshal([]byte(thenEventString), &thenEvents)

			if err != nil {
				fmt.Println(err)
			}
			for variableId, thenEvent := range thenEvents {
				fmt.Println("tuta")
				var token string
				fmt.Println(variableId)
				tools.DBS.QueryRow("SELECT u.token FROM users u INNER JOIN users_variables uv ON uv.userid = u.id WHERE uv.varid = $1", variableId).Scan(&token)
				var label string
				tools.DBS.QueryRow("SELECT label FROM variables WHERE id = $1", variableId).Scan(&label)
				//tools.DBS.Exec("INSERT INTO values (value, varid) VALUES($1, $2)", thenEvent.Value, variableId)
				//tx.Rollback()
				value.SendValueToYandex(0, int(variableId), thenEvent.Value)

				mqtt.Client.Publish(fmt.Sprintf("value/%s/%s", token, label), 0, false, thenEvent.Value)
			}
			fmt.Println(interval.Valid, eventId)
			if !interval.Valid {
				_, err = tools.DBS.Exec("UPDATE events SET isactive = true WHERE id = $1", eventId)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				nowTime = time.Now().Add(time.Minute * time.Duration(int64(interval.Int32)))
				tools.DBS.Exec("UPDATE events SET runtime = $1 WHERE id = $2", nowTime, eventId)
			}
		}
	}
}
