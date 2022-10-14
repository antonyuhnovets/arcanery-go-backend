// Room managing in hub.
// Creation, subscribing, messages and events processing.

package lobby

import (
	"encoding/json"
	"log"

	"github.com/hetonei/arcanery-go-backend/internal/service"
)

type Room struct {
	Id         string                  // room id
	Active     chan bool               // trigger for turning room off/on
	Subs       map[string]Subscription // room subscribers
	Broadcast  chan Message            // redirect messages to subscribers
	Register   chan Subscription       // new subscribe event
	Unregister chan Subscription       // unsubscribe event
}

// Subscriber on room
type Subscription struct {
	Id   string             // subscriber id
	Room string             // room id
	Conn service.Connection // connection interface
}

func CreateRoom(roomId string) *Room {
	room := Room{
		Id:         roomId,
		Subs:       make(map[string]Subscription),
		Active:     make(chan bool),
		Broadcast:  make(chan Message),
		Register:   make(chan Subscription),
		Unregister: make(chan Subscription),
	}
	return &room
}

// Make subscription of given id and connections
func CreateSubscription(id, roomId string, conn service.Connection) Subscription {
	return Subscription{id, roomId, conn}
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

// Start processing events and messages
// Handling according to channel recieved it
func (r *Room) Start() {
	for {
		select {
		case s := <-r.Register: // new subscription request
			log.Println("Room is processing register")
			r.AddSubscriber(s)

		case s := <-r.Unregister: // unsubscribe request
			log.Println("Room is processing unregister")
			r.RemoveSubscriber(s)

		case m := <-r.Broadcast: // new message to broadcast
			log.Println("Room is processing msg")
			r.ProcessMsg(m)

		case trigger := <-r.Active: // close room request
			if !trigger {
				log.Println("Trigger down")

				return
			}
		}
	}
}

// Turn off room, close all channels
// Remove all subscribers
func (r *Room) Shutdown() {
	r.RemoveAllSubscribers()
	r.Active <- false
	r.CloseChannels()
}

func (r *Room) AddSubscriber(subscriber Subscription) {
	if _, ok := r.Subs[subscriber.Id]; !ok {
		r.Subs[subscriber.Id] = subscriber
	}
}

func (r *Room) RemoveSubscriber(subscriber Subscription) {
	delete(r.Subs, subscriber.Id)
}

// Close all room channels
func (r *Room) CloseChannels() {
	close(r.Broadcast)
	close(r.Register)
	close(r.Unregister)
	log.Println("All channels closed")
}

// Clean room from subscribers
func (r *Room) RemoveAllSubscribers() {
	for key, sub := range r.Subs {
		log.Printf("Removing subscriber %s from room", key)
		r.Unregister <- sub
	}
}

// Redirect message to all room subscribers
func (r *Room) ProcessMsg(msg Message) {
	for _, sub := range r.Subs {
		sub.Conn.SendMsg(msg)
		log.Println("Msg was redirected to connections")
	}
}

// Accept message, parse it and redirect to hub
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
}

// Start listening to websocket and message exchange
func (s Subscription) ListenWS(room *Room) {
	go s.Conn.WritePump()
	go s.Conn.ReadPump(s.HandleMsg)
}

func (s Subscription) Unsubscribe(room *Room) {
	room.Unregister <- s

	s.Conn.Close()
}
