package service

type Room struct {
	c ClientService
	g GameService
}

type ClientService interface {
	ConnectToRoom(string)
}

func GetServices(c ClientService, g GameService) Room {
	return Room{
		c: c,
		g: g,
	}
}
