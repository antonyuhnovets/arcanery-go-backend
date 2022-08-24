package room_service

import (
	"log"
	"net/http"

	"github.com/hetonei/arcanery-go-backend/internal/service"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

type ClientWS struct {
	id  string
	Sub *websocket.Subscription
	w   http.ResponseWriter
	req *http.Request
}

func RegisterClient(w http.ResponseWriter, r *http.Request) service.ClientService {
	cws := ClientWS{
		id:  r.RemoteAddr,
		w:   w,
		req: r,
	}

	return cws
}

func GetRoomById(roomId string) *websocket.Room {
	room := websocket.CheckRoomRequest(roomId)
	if room != nil {
		log.Printf("%s in hub", roomId)
		return room
	}

	log.Printf("%s not in hub", roomId)
	return nil
}

func GetClientById(roomId, clientId string) *websocket.Subscription {
	room := GetRoomById(roomId)
	if sub, ok := room.Subs[roomId]; ok {
		return &sub
	}

	return nil
}

func (cws ClientWS) CreateRoom(roomId string) {
	room := websocket.CreateRoom(roomId)
	websocket.H.Register <- room
}

func (cws ClientWS) ConnectToRoom(roomId string) {
	room := GetRoomById(roomId)
	if room == nil {
		log.Println("No rooms with given id found")
		return
	}
	log.Printf("Room %s found", roomId)

	ws, err := websocket.UpgradeToWs(cws.w, cws.req)
	if err != nil {
		log.Fatal(err)
		return
	}

	s := websocket.CreateSubscription(cws.id, roomId, ws)
	cws.Sub = &s

	log.Printf("Substribion %s created", cws.req.RemoteAddr)

	cws.Sub.Subscribe(room)
}

func (cws ClientWS) Disconnect(roomId string) {
	room := GetRoomById(roomId)
	if room == nil {
		log.Fatal("No rooms with given id found")
		return
	}
	sub := GetClientById(cws.req.RemoteAddr, roomId)
	if sub == nil {
		log.Fatal("No client with given id in room")
	}

	room.Unregister <- *sub
}

func (cws ClientWS) DeleteRoom(roomId string) {
	room := GetRoomById(roomId)
	websocket.H.Unregister <- room
	log.Printf("request to unregister room %s from hub sent", roomId)
}
