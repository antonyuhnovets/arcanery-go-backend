package websocket

import "log"

type Room struct {
	Id         string
	Active     chan bool
	Subs       map[string]Subscription
	Broadcast  chan Message
	Register   chan Subscription
	Unregister chan Subscription
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
			r.ProcessMsg(m)
		case trigger := <-r.Active:
			if !trigger {
				log.Println("Trigger down")
				return
			}
		}
	}
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
		close(subscriber.Conn.send)
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
		sub.Conn.send <- msg.Data
		log.Println("Msg was redirected to connections")
	}
}
