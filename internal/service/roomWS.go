package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/hetonei/arcanery-go-backend/pkg/websocket"
)

type RequestRoomWS struct {
	senderId    string
	substribion *websocket.Subscription
	writer      http.ResponseWriter
	request     *http.Request
}

func RegisterRequest(w http.ResponseWriter, r *http.Request) RequestRoomWS {
	rws := RequestRoomWS{
		senderId: r.RemoteAddr[6:],
		writer:   w,
		request:  r,
	}

	return rws
}

func (rws RequestRoomWS) CreateRoom(roomId string) {
	room := websocket.CreateRoom(roomId, EventHandler)
	websocket.H.Register <- room
}

func (rws RequestRoomWS) ConnectToRoom(roomId string) {
	room := GetRoomById(roomId)
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

	s := websocket.CreateSubscription(rws.senderId, roomId, ws)
	rws.substribion = &s

	log.Printf("Substribion %s created", rws.request.RemoteAddr)

	rws.substribion.Subscribe(room)
}

func (rws RequestRoomWS) Disconnect(roomId string) {
	room := GetRoomById(roomId)
	if room == nil {
		log.Fatal("No rooms with given id found")
		return
	}
	sub := GetClientById(rws.request.RemoteAddr, roomId)
	if sub == nil {
		log.Fatal("No client with given id in room")
	}

	room.Unregister <- *sub
}

func (rws RequestRoomWS) DeleteRoom(roomId string) {
	room := GetRoomById(roomId)
	websocket.H.Unregister <- room
	log.Printf("request to unregister room %s from hub sent", roomId)
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

func GetAllSubscribersID(roomId string) []string {
	var subsId []string
	room := GetRoomById(roomId)
	for _, sub := range room.Subs {
		subsId = append(subsId, sub.SubId)
	}

	return subsId
}

func SliceToBytes(slice []string) []byte {
	return []byte(strings.Join(slice, ";\n"))

}

func EventHandler(m websocket.Message) websocket.Message {
	switch m.Event {
	case "connected":
		m.Data = GetAllSubscribersID(m.Room)
	case "chat":
	}

	return m
}
