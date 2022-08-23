package service

type RoomService interface {
	Start()
}

type Room struct {
	id      string
	service RoomService
}

func New(s RoomService, id string) *Room {
	return &Room{
		id:      id,
		service: s,
	}
}

func (r *Room) Start() {
	r.service.Start()
}
