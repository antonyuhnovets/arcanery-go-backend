package service

type RoomService struct {
	c ClientService
	g GameService
}

type ClientService interface {
	ConnectToRoom(string)
	CreateRoom(string)
	Disconnect(string)
}

func GetServices(c ClientService, g GameService) RoomService {
	return RoomService{
		c: c,
		g: g,
	}
}
