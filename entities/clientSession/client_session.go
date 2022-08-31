package clientSession

import (
	"errors"
	"fmt"
	"iot-project/tools"
	"sync"
	"time"
)

type ClientSession struct {
	StartedAt string `json:"startedAt"`
	UserAgent string `json:"userAgent"`
	SessionID string `json:"userAgent"`
	Channel   chan string
}
type СlientSessions struct {
	sync.Mutex
	Sessions []ClientSession
}

type clientSessionControllerImpl struct {
	storage sync.Map
}

type clientSessionController interface {
	AddNewSessionToStorage(actorId uint32, channel chan string, userAgent string) (string, error)
	RemoveSessionByChannel(actorId uint32, channel chan string) error
	RemoveSessionsByActorId(actorId uint32)
	RemoveSessionById(actorId uint32, id string) error
	GetSessionsByActorId(actorId uint32) (*[]ClientSession, error)
	GetSessionsBySessionId(actorId uint32, sessionId string) (*ClientSession, error)
	GetAllSessions() sync.Map
	CountSessionsByActorID(actorId uint32) int
}

func (clientSessionStorage *clientSessionControllerImpl) AddNewSessionToStorage(actorId uint32, channel chan string, userAgent string) (string, error) {
	sessionsByActorId, _ := clientSessionStorage.storage.LoadOrStore(actorId, &СlientSessions{Sessions: make([]ClientSession, 0, tools.MaxSessionsFromClient)})
	sessions, ok := sessionsByActorId.(*СlientSessions)

	if !ok {
		panic("clientSessionStorage must be contain only pointers on clientSessions struct")
	}

	sessions.Lock()

	defer sessions.Unlock()

	for i := 0; i < len(sessions.Sessions); i++ {
		if sessions.Sessions[i].Channel == channel {
			return "", errors.New("this channel has already been added")
		}
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	currentTime := time.Now().In(loc)
	sessionId := fmt.Sprint(time.Now().UnixNano())

	if len(sessions.Sessions) == cap(sessions.Sessions) {
		//remove first session, move slice elements
		sessions.Sessions[0].Channel <- fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Error")
		for i := 0; i <= tools.MaxSessionsFromClient-1; i++ {
			if i+1 < len(sessions.Sessions) {
				sessions.Sessions[i] = sessions.Sessions[i+1]
			}
		}
		sessions.Sessions[tools.MaxSessionsFromClient-1] = ClientSession{Channel: channel, StartedAt: currentTime.Format(time.RFC3339), UserAgent: userAgent, SessionID: sessionId}
	} else {
		sessions.Sessions = append(sessions.Sessions, ClientSession{Channel: channel, StartedAt: currentTime.Format(time.RFC3339), UserAgent: userAgent, SessionID: sessionId})
	}
	return sessionId, nil
}

func (clientSessionStorage *clientSessionControllerImpl) RemoveSessionByChannel(actorId uint32, channel chan string) error {
	sessionsByActorId, ok := clientSessionStorage.storage.Load(actorId)

	if !ok {
		return errors.New("no such active actor")
	}

	sessions, ok := sessionsByActorId.(*СlientSessions)

	if !ok {
		panic("clientSessionStorage must be contain only pointers on clientSessions struct")
	}

	sessions.Lock()
	defer sessions.Unlock()

	for i := 0; i < len(sessions.Sessions); i++ {
		if sessions.Sessions[i].Channel == channel {
			sessions.Sessions[i].Channel <- fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Error")
			sessions.Sessions[i] = sessions.Sessions[len(sessions.Sessions)-1]
			sessions.Sessions = sessions.Sessions[:len(sessions.Sessions)-1]
			return nil
		}
	}
	return errors.New("no such channel for this actor")
}

func (clientSessionStorage *clientSessionControllerImpl) RemoveSessionById(actorId uint32, sessionId string) error {
	sessionsByActorId, ok := clientSessionStorage.storage.Load(actorId)

	if !ok {
		return errors.New("no such active actor")
	}

	sessions, ok := sessionsByActorId.(*СlientSessions)

	if !ok {
		panic("clientSessionStorage must be contain only pointers on clientSessions struct")
	}

	sessions.Lock()
	defer sessions.Unlock()

	for i := 0; i < len(sessions.Sessions); i++ {
		if sessions.Sessions[i].SessionID == sessionId {
			sessions.Sessions[i].Channel <- fmt.Sprintf("{\"errno\":%d,\"message\":\"%s\"}", 1, "Error")
			sessions.Sessions[i] = sessions.Sessions[len(sessions.Sessions)-1]
			sessions.Sessions = sessions.Sessions[:len(sessions.Sessions)-1]
			return nil
		}
	}
	return errors.New("no such sessionID for this actor")
}

func (clientSessionStorage *clientSessionControllerImpl) RemoveSessionsByActorId(actorId uint32) {
	clientSessionStorage.storage.Delete(actorId)
}

func (clientSessionStorage *clientSessionControllerImpl) GetSessionsByActorId(actorId uint32) (*[]ClientSession, error) {
	sessionsByActorId, ok := clientSessionStorage.storage.Load(actorId)

	if !ok {
		return nil, errors.New("no such active actor")
	}
	sessions := sessionsByActorId.(*СlientSessions)
	return &sessions.Sessions, nil
}

func (clientSessionStorage *clientSessionControllerImpl) GetSessionsBySessionId(actorId uint32, sessionId string) (*ClientSession, error) {
	sessionsByActorId, ok := clientSessionStorage.storage.Load(actorId)

	if !ok {
		return nil, errors.New("no such active actor")
	}
	sessions := sessionsByActorId.(*СlientSessions)
	if !ok {
		panic("clientSessionStorage must be contain only pointers on clientSessions struct")
	}

	for i := 0; i < len(sessions.Sessions); i++ {
		if sessions.Sessions[i].SessionID == sessionId {
			return &sessions.Sessions[i], nil
		}
	}
	return nil, errors.New("no such sessionID for this actor")
}

func (clientSessionStorage *clientSessionControllerImpl) CountSessionsByActorID(actorId uint32) int {

	sessionsByActorId, ok := clientSessionStorage.storage.Load(actorId)
	if !ok {
		return 0
	}
	sessions := sessionsByActorId.(*СlientSessions)
	return len(sessions.Sessions)
}

func (clientSessionStorage *clientSessionControllerImpl) GetAllSessions() sync.Map {
	return clientSessionStorage.storage
}

var ClientSessionController clientSessionController = &clientSessionControllerImpl{}
