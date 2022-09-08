package service

import "context"

type Connection interface {
	WritePump()
	SendMsg(interface{})
	ReadPump(func([]byte))
	SetConnectionId(string)
	GetConnectionId() string
	Close()
}

type Services struct {
	Ctx  context.Context
	Room RoomService
	Repo RepositoryService
}

type RoomService interface {
	CreateRoom(string)
	ConnectToRoom(string)
	Disconnect(string)
	DeleteRoom(string)
}

type RepositoryService interface {
	Create(interface{}) error
	ReadById(interface{}, int64) error
	ReadAll(interface{}) error
	UpdateById(interface{}, int64) error
	DeleteById(int64) error
	DeleteAll() (int64, error)
}

type ConnectionDB interface {
	GetRepoService(string) RepositoryService
}

func (s *Services) SetRoomService(srv RoomService) {
	s.Room = srv
}

func (s *Services) SetRepoService(srv RepositoryService) {
	s.Repo = srv
}
