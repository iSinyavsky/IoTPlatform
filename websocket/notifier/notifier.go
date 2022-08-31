package notifier

import (
	"fmt"
	"github.com/gorilla/websocket"
	"iot-project/entities/clientSession"
	"strings"
	"time"
)

//Notifier is mechanism which send updates for all active sessions of actor via websocket

const closeConnectionMessage = "CLOSE_CONNECTION"

func StartListenNotifyForActor(actorId uint32, connection *websocket.Conn, channel chan string) {
	//Starting goroutine for handle close event from client
	go func() {
		for {
			messageType, _, _ := connection.ReadMessage()
			if messageType == -1 { //if connection is closed
				channel <- closeConnectionMessage
				return
			}
		}
	}()

	//If connection is active
	for {
		select {
		case message := <-channel:
			if message == closeConnectionMessage {
				fmt.Println("CLOSE CONNECTION FROM closeConnectionMessage")
				connection.Close()
				clientSession.ClientSessionController.RemoveSessionByChannel(actorId, channel)
				break
			}
			err := connection.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				connection.Close()
				err = clientSession.ClientSessionController.RemoveSessionByChannel(actorId, channel)
				break
			}
			continue
		}
		break
	}

	/*
		If connection is closed.
		The channel does not close immediately so that there is no panic in the event that a update has come
		and the channel has not yet been removed from the storage
	*/
	for {
		select {
		case <-channel: //connection is closed and we will ignore updates
			continue
		case <-time.After(10 * time.Second): //close channel after 10 seconds
			close(channel)
			return
		}
	}
}

//func NotifyByVertexId(vertexId uint64, update string) error {
//	vertex := s3.GetSuperNodeByVertexId(vertexId)
//
//	vertexActors := vertex.CSet.Actors
//
//	for actorId, _ := range vertexActors{
//		err := NotifyByActorId(actorId, update)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func NotifyByVertices(vertices []uint64, update string) error {
//	var clustersActors = make([]uint32, 0, 100)
//
//	for i := 0; i < len(vertices); i++ {
//
//
//		vertex := s3.GetSuperNodeByVertexId(vertices[i])
//
//		CSet := vertex.CSet.Actors
//		vertexActors := make([]uint32, 0, 32)
//		for actorId, _ := range CSet{
//			vertexActors = append(vertexActors, actorId)
//		}
//
//		clustersActors = append(clustersActors, vertexActors...)
//	}
//
//	clustersActors = uniq.Uint32(clustersActors)
//
//	for i := 0; i < len(clustersActors); i++ {
//		if clustersActors[i] == 0 {
//			break
//		}
//		err := NotifyByActorId(clustersActors[i], update)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func NotifyByActorIDToOtherSessions(actorId uint32, sessionId string, update string) error {
	listeners, err := clientSession.ClientSessionController.GetSessionsByActorId(actorId)
	if err != nil {
		if err.Error() == "no such active actor" {
			return nil
		}
		return err
	}
	for j := 0; j < len(*listeners); j++ {
		fmt.Println("session", (*listeners)[j].SessionID)
		if (*listeners)[j].SessionID != sessionId {
			(*listeners)[j].Channel <- update
		}

	}

	return nil
}

type notifyTxImpl struct {
	updates map[uint32]*strings.Builder
}
type notifyTx interface {
	Commit() error
	AddNotifyByActorId(actorId uint32, update string) error
}

func Begin() notifyTx {
	return &notifyTxImpl{updates: make(map[uint32]*strings.Builder)}
}

func (tx notifyTxImpl) AddNotifyByActorId(actorId uint32, update string) error {
	if _, ok := tx.updates[actorId]; !ok {
		tx.updates[actorId] = &strings.Builder{}
	}
	tx.updates[actorId].WriteString(update)
	tx.updates[actorId].WriteString(",")
	return nil
}

func (tx *notifyTxImpl) Commit() error {
	for aID, actorUpdates := range tx.updates {
		actorUpdates.WriteString("\"_0\":{}")
		update := fmt.Sprintf("{\"data\":{\"vertices\": {%s}}}", actorUpdates.String())
		err := NotifyByActorId(aID, update)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func NotifyByActorId(actorId uint32, update string) error {
	listeners, err := clientSession.ClientSessionController.GetSessionsByActorId(actorId)
	if err != nil {
		if err.Error() == "no such active actor" {
			return nil
		}
		return err
	}

	for j := 0; j < len(*listeners); j++ {
		(*listeners)[j].Channel <- update
	}

	return nil
}
