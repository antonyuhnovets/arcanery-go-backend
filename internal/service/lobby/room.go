package lobby

import (
	"encoding/json"
	"log"

	"github.com/hetonei/arcanery-go-backend/internal/service"
)

type Subscription struct {
	SubId string
	Room  string
	Conn  service.Connection
}

func CreateSubscription(id, roomId string, conn service.Connection) Subscription {
	return Subscription{id, roomId, conn}
}

func (s Subscription) HandleMsg(msg []byte) {
	incMsg := &Message{}
	err := json.Unmarshal(msg, incMsg)

	incMsg.Room = s.Room
	if err != nil {
		log.Println(err)
	}

	H.Broadcast <- *incMsg
}

func (s Subscription) Subscribe(room *Room) {
	room.Register <- s

	go s.Conn.WritePump()
	go s.Conn.ReadPump(s.HandleMsg)
}

type Room struct {
	Id           string
	Active       chan bool
	Subs         map[string]Subscription
	Broadcast    chan Message
	Register     chan Subscription
	Unregister   chan Subscription
	EventHandler func(Message) Message
}

func CreateRoom(roomId string, handler func(Message) Message) *Room {
	room := Room{
		Id:           roomId,
		Subs:         make(map[string]Subscription),
		Active:       make(chan bool),
		Broadcast:    make(chan Message),
		Register:     make(chan Subscription),
		Unregister:   make(chan Subscription),
		EventHandler: handler,
	}
	return &room
}

func (r *Room) Start() {
	for {
		select {
		case s := <-r.Register:
			log.Println("Room is processing register")
			r.AddSubscriber(s)
		case s := <-r.Unregister:
			log.Println("Room is processing unregister")
			r.RemoveSubscriber(s)
		case m := <-r.Broadcast:
			log.Println("Room is processing msg")
			msg := r.EventHandler(m)
			r.ProcessMsg(msg)
		case trigger := <-r.Active:
			if !trigger {
				log.Println("Trigger down")
				return
			}
		}
	}
}

func RegisterRoom(r *Room) {
	H.Register <- r
}

func Unregister(r *Room) {
	H.Unregister <- r
}

func (r *Room) Shutdown() {
	r.RemoveAllSubscribers()
	r.Active <- false
	r.CloseChannels()
}

func (r *Room) AddSubscriber(subscriber Subscription) {
	if _, ok := r.Subs[subscriber.SubId]; !ok {
		r.Subs[subscriber.SubId] = subscriber
	}
}

func (r *Room) RemoveSubscriber(subscriber Subscription) {
	if _, ok := r.Subs[subscriber.SubId]; ok {
		subscriber.Conn.Close()
		delete(r.Subs, subscriber.SubId)
	}
}

func (r *Room) CloseChannels() {
	close(r.Broadcast)
	close(r.Register)
	close(r.Unregister)
	log.Println("All channels closed")
}

func (r *Room) RemoveAllSubscribers() {
	for key, sub := range r.Subs {
		log.Printf("Removing subscriber %s from room", key)
		r.Unregister <- sub
	}
}

func (r *Room) ProcessMsg(msg Message) {
	for _, sub := range r.Subs {
		sub.Conn.SendMsg(msg)
		log.Println("Msg was redirected to connections")
	}
}

func GetRoomById(roomId string) *Room {
	room := CheckRoomRequest(roomId)
	if room != nil {
		log.Printf("%s in hub", roomId)
		return room
	}

	log.Printf("%s not in hub", roomId)
	return nil
}

func GetClientById(roomId, clientId string) *Subscription {
	room := GetRoomById(roomId)
	if sub, ok := room.Subs[roomId]; ok {
		return &sub
	}

	return nil
}
