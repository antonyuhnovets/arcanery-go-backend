// Room managing with websocket.

package ws

import (
	"log"
	"net/http"

	"github.com/hetonei/arcanery-go-backend/internal/service/lobby"
	"github.com/hetonei/arcanery-go-backend/pkg/uuid"
	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

// Request on websocket room
type RequestWS struct {
	writer  http.ResponseWriter
	request *http.Request
}

// Get http request and add writer for ws
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

// Connect to room by id, upgrade connection protocol to websocket
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

	// Make subscription and send on register
	conn := websocket.MakeConnection(ws)
	newId := uuid.GenerateId()
	sub := lobby.CreateSubscription(newId, roomId, conn)

	log.Printf("Subscribion %s created", newId)

	sub.Subscribe(room)
	sub.ListenWS(room)

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

// Send room unregister request to hub
func (rws RequestWS) DeleteRoom(roomId string) {
	room := lobby.GetRoomById(roomId)
	lobby.Unregister(room)

	log.Printf("request to unregister room %s from hub sent", roomId)
}
