package service

import "context"

type Lobby struct {
	id string
	// users map[string]User
}

type Services struct {
	Ctx  context.Context
	Room RoomService
}

type RoomService interface {
	CreateRoom(string)
	ConnectToRoom(string)
	Disconnect(string)
	DeleteRoom(string)
}

func (s *Services) SetRoomService(srv RoomService) {
	s.Room = srv
}

func NewLobby(creator Services, roomId string) *Lobby {
	creator.Room.CreateRoom(roomId)
	creator.Room.ConnectToRoom(roomId)
	return &Lobby{
		id: roomId,
		// users: map[string]Services{creator.Id: creator},
	}
}

type Connection interface {
	WritePump()
	SendMsg(interface{})
	ReadPump(func([]byte))
	SetConnectionId(string)
	GetConnectionId() string
	Close()
}

type PlayerService interface {
	ChoseCards(...string)
	PlayCard(string)
	GetReady()
	ChangeCard(string, string)
}

type GameService interface {
	StartGame()
	EndTurn()
	EndRound()
	EndGame()
}
