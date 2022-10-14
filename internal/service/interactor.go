// Represents abstract interfaces with buisness logic.

package service

import "context"

// Interface for tunnel connetion (like websocket)
// Controls msg exchange process (handshake)
type Connection interface {
	WritePump()
	SendMsg(interface{})
	ReadPump(func([]byte))
	SetConnectionId(string)
	GetConnectionId() string
	Close()
}

// Interface for room manage
type RoomService interface {
	CreateRoom(string)
	ConnectToRoom(string)
	Disconnect(string)
	DeleteRoom(string)
}

// Abstract CRUD interface
// Saving, processing data in repository
type RepositoryService interface {
	Create(interface{}) error
	ReadById(interface{}, int64) error
	ReadAll(interface{}) error
	UpdateById(interface{}, int64) error
	DeleteById(int64) error
	DeleteAll() (int64, error)
}

type ConnectionRepo interface {
	GetRepoService(string) RepositoryService
}

// Essence with main application services
// Services is abstract interfaces with buisness logic
type Services struct {
	Ctx  context.Context
	Room RoomService
	Repo RepositoryService
}

func (s *Services) SetRoomService(srv RoomService) {
	s.Room = srv
}

func (s *Services) SetRepoService(srv RepositoryService) {
	s.Repo = srv
}
