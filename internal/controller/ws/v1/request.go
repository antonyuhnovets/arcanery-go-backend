package ws

import (
	"log"
	"net/http"

	"github.com/hetonei/arcanery-go-backend/internal/service/lobby"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

type RequestWS struct {
	senderId    string
	substribion lobby.Subscription
	writer      http.ResponseWriter
	request     *http.Request
}

func RegisterRequest(w http.ResponseWriter, r *http.Request) RequestWS {
	rws := RequestWS{
		senderId: r.RemoteAddr[6:],
		writer:   w,
		request:  r,
	}

	return rws
}

func (rws RequestWS) CreateRoom(roomId string) {
	room := lobby.CreateRoom(roomId, lobby.EventHandler)
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
	s := lobby.CreateSubscription(rws.senderId, roomId, conn)
	rws.substribion = s

	log.Printf("Substribion %s created", rws.request.RemoteAddr)

	rws.substribion.Subscribe(room)
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

	room.Unregister <- *sub
}

func (rws RequestWS) DeleteRoom(roomId string) {
	room := lobby.GetRoomById(roomId)
	lobby.Unregister(room)
	log.Printf("request to unregister room %s from hub sent", roomId)
}

// func ServeWs(w http.ResponseWriter, r *http.Request, roomId string) {
// 	room := CheckRoomRequest(roomId)

// 	ws, err := UpgradeToWs(w, r)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	c := MakeConnection(ws)
// 	s := CreateSubscription(r.RemoteAddr, roomId, c)

// 	s.Subscribe(room)
// }
