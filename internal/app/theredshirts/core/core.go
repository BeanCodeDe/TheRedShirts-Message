package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/adapter"
	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/db"
	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/util"
	"github.com/google/uuid"
)

type (

	//Facade
	CoreFacade struct {
		db            db.DB
		lobbyAdapter  *adapter.LobbyAdapter
		lobbyPlayerId uuid.UUID
	}

	Core interface {
		//Message
		CreateMessage(context *util.Context, message *Message) error
		GetMessages(context *util.Context, playerId uuid.UUID, lobbyId uuid.UUID, number int) ([]*Message, error)
	}

	//Objects
	Message struct {
		ID       uuid.UUID
		SendTime time.Time
		LobbyId  uuid.UUID
		PlayerId uuid.UUID
		Number   int
		Topic    string
		Message  map[string]interface{}
	}

	Player struct {
		ID          uuid.UUID
		LobbyId     uuid.UUID
		LastRefresh time.Time
	}
)

var (
	ErrWrongLobbyPassword = errors.New("wrong password")
)

func NewCore() (Core, error) {
	db, err := db.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("error while initializing database: %v", err)
	}
	lobbyAdapter := adapter.NewLobbyAdapter()
	lobbyPlayerId, err := util.GetEnvUUID("LOBBY_USER")
	if err != nil {
		return nil, fmt.Errorf("error while loading lobby user env: %v", err)
	}
	core := &CoreFacade{db: db, lobbyAdapter: lobbyAdapter, lobbyPlayerId: lobbyPlayerId}
	core.startCleanUp()
	return core, nil
}
