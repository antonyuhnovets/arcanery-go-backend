// Room managing in hub.
// Creation, subscribing, messages and events processing.

package lobby

import (
	"encoding/json"
	"log"
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
	Id     string // subscriber id
	RoomId string // room id
	Room   *Room
	Conn   Connection // connection interface
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
func (r *Room) CreateSubscription(id string, conn interface{}) *Subscription {
	i, ok := conn.(Connection)
	if !ok {
		log.Println("No connection to subscribe")
		return nil
	}
	return &Subscription{
		Id:     id,
		RoomId: r.Id,
		Room:   r,
		Conn:   i,
	}
}

func (r *Room) GetClientById(clientId string) *Subscription {
	if sub, ok := r.Subs[r.Id]; ok {
		return &sub
	}
	log.Println("Subscriber not found")
	return nil
}

// Start processing events and messages
// Handling according to channel recieved it
func (r *Room) Start() {
	for {
		select {
		case s := <-r.Register: // new subscription request
			log.Printf("Room %s is processing register sub %s", s.RoomId, s.Id)
			r.AddSubscriber(s)

		case s := <-r.Unregister: // unsubscribe request
			log.Printf("Room %s is processing unregister sub %s", s.RoomId, s.Id)
			r.RemoveSubscriber(s)

		case m := <-r.Broadcast: // new message to broadcast
			log.Printf("Room %s is processing msg", m.Room)
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
		log.Printf("Subscriber %s registred", subscriber.Id)
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

func (r *Room) GetAllSubscribers() map[string]interface{} {
	l := make(map[string]interface{})
	for key, sub := range r.Subs {
		l[key] = sub
	}

	return l
}

// Redirect message to all room subscribers
func (r *Room) ProcessMsg(msg Message) {
	for _, sub := range r.Subs {
		sub.Conn.SendMsg(msg)
	}
	log.Println("Msg was redirected to connections")
}

func (r *Room) GetRoomInfo() *Room {
	return r
}

// Accept message, parse it and redirect to hub
func (s Subscription) HandleMsg(msg []byte) {
	incMsg := &Message{}
	err := json.Unmarshal(msg, incMsg)

	incMsg.Room = s.RoomId
	if err != nil {
		log.Println(err)
	}

	s.Room.Broadcast <- *incMsg
}

func (s Subscription) Subscribe(r *Room) {
	r.Register <- s

	s.Conn.Listen(s.HandleMsg)
}

func (s Subscription) Unsubscribe(room *Room) {
	room.Unregister <- s

	s.Conn.Close()
}
