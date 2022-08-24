package service

type ClientService interface {
	CreateRoom(string)
	DeleteRoom(string)
	ConnectToRoom(string)
	Disconnect(string)
}
