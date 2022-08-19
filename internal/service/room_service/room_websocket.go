package websocket

import (
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

type RoomService struct {
	l *websocket.Lobby
}

func (r *RoomService) Start() {
	go r.l.Run()
}
