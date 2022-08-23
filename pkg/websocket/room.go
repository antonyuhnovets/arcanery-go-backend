package websocket

type Room struct {
	Id         string
	active     bool
	subs       map[string]Subscription
	Broadcast  chan Message
	Register   chan Subscription
	Unregister chan Subscription
}

func CreateRoom(roomId string) *Room {
	room := Room{
		Id:         roomId,
		active:     true,
		subs:       make(map[string]Subscription),
		Broadcast:  make(chan Message),
		Register:   make(chan Subscription),
		Unregister: make(chan Subscription),
	}
	return &room
}

func (r *Room) Start() {
	select {
	case s := <-r.Register:
		r.AddSubscriber(s)
	case s := <-r.Unregister:
		r.RemoveSubscriber(s)
	case m := <-r.Broadcast:
		r.ProcessMsg(m)
	}
}

func (r *Room) Shutdown() {
	close(r.Broadcast)
	close(r.Register)
	close(r.Unregister)
}

func (r *Room) AddSubscriber(subscriber Subscription) {
	if _, ok := r.subs[subscriber.subId]; !ok {
		r.subs[subscriber.subId] = subscriber
	}
}

func (r *Room) RemoveSubscriber(subscriber Subscription) {
	if _, ok := r.subs[subscriber.subId]; ok {
		close(subscriber.Conn.send)
		delete(r.subs, subscriber.subId)
	}
}

func (r *Room) ProcessMsg(msg Message) {
	for _, sub := range r.subs {
		sub.Conn.send <- msg.Data
	}
}
