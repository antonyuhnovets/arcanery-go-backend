// Package represent lobby realisation. Main logic of hub with rooms, message and events processing.
// File implement event handler.

package lobby

import (
	"encoding/json"
	"errors"
	"log"
)

type Request struct {
	Event  string
	Params map[string]string
	Data   map[string]interface{}
}

func NewEventHub(roomHub *Hub, eventMap map[string]string) *EventHub {
	h := &EventHub{
		roomHub:      roomHub,
		events:       eventMap,
		OnCreate:     make(chan *Request),
		OnConnect:    make(chan *Request),
		OnDelete:     make(chan *Request),
		OnDisconnect: make(chan *Request),
		OnMessage:    make(chan *Request),
		OnErr:        make(chan error),
	}

	// h.SetEvents(eventMap)

	return h
}

type EventHub struct {
	roomHub      *Hub
	events       map[string]string
	OnCreate     chan *Request
	OnConnect    chan *Request
	OnDelete     chan *Request
	OnDisconnect chan *Request
	OnMessage    chan *Request
	OnErr        chan error
}

func (eh *EventHub) Run() {

	for {
		select {
		case r := <-eh.OnCreate:
			eh.Create(r)

		case r := <-eh.OnConnect:
			eh.Connect(r)

		case r := <-eh.OnDelete:
			eh.Delete(r)

		case r := <-eh.OnDisconnect:
			eh.Disconnect(r)

		case r := <-eh.OnMessage:
			eh.Message(r)

		case er := <-eh.OnErr:
			log.Println(er)
		}
	}
}

func (eh *EventHub) MakeRequest(r *Request) {
	log.Printf("Sending %s request on EventHub", r.Event)

	switch eh.events[r.Event] {
	case "OnCreate":
		eh.OnCreate <- r
	case "OnConnect":
		eh.OnConnect <- r
	case "OnDelete":
		eh.OnDelete <- r
	case "OnDisconnect":
		eh.OnDisconnect <- r
	case "OnMessage":
		eh.OnMessage <- r

		log.Printf("Request %s sended", r.Event)
	}
}

func (eh *EventHub) Create(r *Request) {
	roomId, ok := r.Params["room"]
	if !ok {
		eh.OnErr <- errors.New("OnCreate failed: no roomId param")
		return
	}

	if eh.roomHub.CheckRoomInHub(roomId) {
		eh.OnErr <- errors.New("OnCreate failed: exiting room")
		return
	}

	room := CreateRoom(roomId)
	eh.roomHub.RegisterRoom(room)
}

func (eh *EventHub) Delete(r *Request) {
	roomId, ok := r.Params["room"]
	if !ok {
		eh.OnErr <- errors.New("OnDelete failed: no roomId param")
		return
	}

	room := eh.roomHub.GetRoomById(roomId)
	if room == nil {
		eh.OnErr <- errors.New("OnDelete failed: room not found")
		return
	}

	eh.roomHub.UnregisterRoom(room)
}

func (eh *EventHub) Connect(r *Request) {
	roomId, ok := r.Params["room"]
	if !ok {
		eh.OnErr <- errors.New("OnConnect failed: no roomId param")
		return
	}

	room := eh.roomHub.GetRoomById(roomId)
	if room == nil {
		eh.OnErr <- errors.New("OnConnect failed: room not found")
		return
	}

	c := NewConnectionManager(r.Params["connType"], r.Data)
	sub := room.CreateSubscription(r.Params["connId"], c)
	sub.Subscribe(room)
}

func (eh *EventHub) Disconnect(r *Request) {
	roomId, ok := r.Params["room"]
	if !ok {
		eh.OnErr <- errors.New("OnDisconnect failed: no roomId param")
		return
	}

	room := eh.roomHub.GetRoomById(roomId)
	if room == nil {
		eh.OnErr <- errors.New("OnDisconnect failed: room not found")
		return
	}

	sub := room.GetClientById(r.Params["connId"])
	if sub == nil {
		eh.OnErr <- errors.New("OnDisconnect fail: no client in group")
		return
	}

	(*sub).Unsubscribe(room)
}

func (eh *EventHub) Message(r *Request) {
	roomId, ok := r.Params["room"]
	if !ok {
		eh.OnErr <- errors.New("OnAction failed: no roomId param")
		return
	}

	room := eh.roomHub.GetRoomById(roomId)
	if room == nil {
		eh.OnErr <- errors.New("OnAction failed: room not found")
		return
	}

	eh.roomHub.Broadcast <- Message{
		Room:  roomId,
		Event: r.Event,
		Data:  r.Data,
	}
}

func (eh *EventHub) RoomCheck(id string) error {
	r := eh.roomHub.CheckRoomRequest(id)
	if r == nil {
		return errors.New("no requested room")
	}
	return nil
}

// Pull some information from message
func (r *Request) PullData(data string) interface{} {
	b, err := json.Marshal(r.Data)
	if err != nil {
		log.Println(err)
	}

	value := GetDataObj(b, data)

	return value
}

// Find value of needed field
func GetDataObj(data []byte, key string) interface{} {
	d := make(map[string]interface{})
	err := json.Unmarshal(data, &d)

	if err != nil {
		log.Println(err)
	}

	return d[key]
}

// Get list of all subscribers ids in room
func GetAllSubscribersID(roomId string) []string {
	var subsId []string
	room := H.GetRoomById(roomId)
	for _, sub := range room.Subs {
		subsId = append(subsId, sub.Id)
	}

	return subsId
}
