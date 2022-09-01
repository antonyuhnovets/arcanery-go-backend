package websocket

import "log"

type Message struct {
	Room  string      `json:"rid"`
	Event string      `json:"type"`
	Data  interface{} `json:"data"`
}

type Hub struct {
	// Registered connections.
	Rooms map[string]*Room

	// Inbound messages from the connections.
	Broadcast chan Message

	// Register requests from the connections.
	Register chan *Room

	// Unregister requests from connections.
	Unregister chan *Room
}

var H = &Hub{
	Broadcast:  make(chan Message),
	Register:   make(chan *Room),
	Unregister: make(chan *Room),
	Rooms:      make(map[string]*Room),
}

func (h *Hub) Run() {
	for {
		select {
		case r := <-h.Register:
			log.Println("Hub is processing register")
			h.AddRoom(r)
		case r := <-h.Unregister:
			log.Println("Hub is processing unregister")
			h.RemoveRoom(r)
		case m := <-h.Broadcast:
			log.Println("Hub is processing msg")
			h.RedirectMsg(m)
		}
	}
}

func (h *Hub) AddRoom(room *Room) {
	h.Rooms[room.Id] = room
	log.Printf("%s added in hub", room.Id)
	go room.Start()
	log.Println("The room started")
}

func (h *Hub) RemoveRoom(room *Room) {
	room.Shutdown()
	log.Println("Shutdown correct, deleting from hub")
	delete(h.Rooms, room.Id)
	log.Printf("room has been removed from hub")
}

func (h *Hub) RedirectMsg(msg Message) {
	h.Rooms[msg.Room].Broadcast <- msg
	log.Println("Msg redirected")
}

func (h *Hub) CheckRoomInHub(roomId string) bool {
	if _, ok := h.Rooms[roomId]; ok {
		log.Println(true)
		return true
	}
	return false
}
