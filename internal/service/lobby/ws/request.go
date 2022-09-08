package ws

import (
	"log"
	"net/http"

	// "github.com/hetonei/arcanery-go-backend/internal/service"
	"github.com/hetonei/arcanery-go-backend/internal/service/lobby"
	"github.com/hetonei/arcanery-go-backend/pkg/uuid"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

type RequestWS struct {
	writer  http.ResponseWriter
	request *http.Request
}

func RegisterRequest(w http.ResponseWriter, r *http.Request) RequestWS {
	rws := RequestWS{
		writer:  w,
		request: r,
	}

	return rws
}

func (rws RequestWS) CreateRoom(roomId string) {
	room := lobby.CreateRoom(roomId)
	lobby.RegisterRoom(room)
}

func (rws RequestWS) ConnectToRoom(roomId string) {
	room := lobby.GetRoomById(roomId)
	if room == nil {
		log.Println("No rooms with given id found")
		return
	}
	log.Printf("Room %s found", roomId)

	ws, err := websocket.UpgradeToWs(rws.writer, rws.request)
	if err != nil {
		log.Fatal(err)
		return
	}

	conn := websocket.MakeConnection(ws)
	newId := uuid.GenerateId()
	s := lobby.CreateSubscription(newId, roomId, conn)

	log.Printf("Subscribion %s created", newId)

	s.Subscribe(s.SubId, room)
	s.ListenWS(room)

}

func (rws RequestWS) Disconnect(roomId string) {
	room := lobby.GetRoomById(roomId)
	if room == nil {
		log.Fatal("No rooms with given id found")
		return
	}

	sub := lobby.GetClientById(rws.request.RemoteAddr, roomId)
	if sub == nil {
		log.Fatal("No client with given id in room")
	}

	// Maybe Error dereference
	sub.Unsubscribe(room)
}

func (rws RequestWS) DeleteRoom(roomId string) {
	room := lobby.GetRoomById(roomId)
	lobby.Unregister(room)
	log.Printf("request to unregister room %s from hub sent", roomId)
}
