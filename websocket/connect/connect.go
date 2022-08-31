package connect

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"iot-project/entities/clientSession"
	"iot-project/tools/auth"
	"iot-project/websocket/notifier"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    128,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	http.HandleFunc("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if _, ok := err.(websocket.HandshakeError); ok {
			//if connection is not websocket
			w.Write([]byte(fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Websocket Error")))
			return
		} else if err != nil {
			w.Write([]byte(fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, err)))
			return
		}

		token, err := auth.GetSessionToken(r)
		if err != nil {
			fmt.Println(err)
			connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Bad token")))
			//w.Write([]byte(tools.ExtResult(err)))
			return
		}

		actorId, err := auth.GetActorIdByAuthToken(token)
		if err != nil {
			fmt.Println(err)
			connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Bad token")))
			connection.Close()
			return
		}

		if actorId == 0 {
			connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Bad token")))
			connection.Close()
			return
		}

		fmt.Println(fmt.Sprintf("connect actor %d", actorId))
		var channel = make(chan string, 50)
		userAgent := r.Header.Get("User-Agent")
		_, err = clientSession.ClientSessionController.AddNewSessionToStorage(actorId, channel, userAgent)
		if err != nil {
			connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Error")))
			connection.Close()
			return
		}

		//send ok and start notifier to send updates
		connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"errno\":%d}", 0)))

		go notifier.StartListenNotifyForActor(actorId, connection, channel)
	})
}
