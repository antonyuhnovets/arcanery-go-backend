// Implementation of hub, that managing (register, unregister) rooms.
// Hub process all requests, events, messages before they redirected to rooms.

package lobby

import "log"

type Message struct {
	Room  string      `json:"rid"`  // room id of message
	Event string      `json:"type"` // type of message (or event)
	Data  interface{} `json:"data"` // message content
}

type Hub struct {
	Rooms      map[string]*Room // rooms with id
	Broadcast  chan Message     // broadcast messages to rooms
	Register   chan *Room       // register new room event
	Unregister chan *Room       // remove room event
}

var H = &Hub{
	Broadcast:  make(chan Message),
	Register:   make(chan *Room),
	Unregister: make(chan *Room),
	Rooms:      make(map[string]*Room),
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		Register:   make(chan *Room),
		Unregister: make(chan *Room),
		Rooms:      make(map[string]*Room),
	}
}

// Start process events by hub
func (h *Hub) Run() {
	for {
		select {
		case r := <-h.Register: // new room
			log.Println("Hub is processing register")
			h.AddRoom(r)

		case r := <-h.Unregister: // delete room
			log.Println("Hub is processing unregister")
			h.RemoveRoom(r)

		case msg := <-h.Broadcast: // incoming message/event
			// (msg)HandleEvent()
			h.RedirectMsg(msg)

		}
	}
}

// Send room registration request to hub
func (h *Hub) RegisterRoom(r *Room) {
	h.Register <- r
}

// Send room unregister request to hub
func (h *Hub) UnregisterRoom(r *Room) {
	h.Unregister <- r
}

func (h *Hub) SendMessage(msg Message) {
	h.Broadcast <- msg
}

// Add room to hub, run gorutine
func (h *Hub) AddRoom(room *Room) {
	h.Rooms[room.Id] = room
	log.Printf("%s added in hub", room.Id)

	go room.Start()
	log.Println("The room started")
}

// Remove room from hub
func (h *Hub) RemoveRoom(room *Room) {
	room.Shutdown()
	log.Println("Shutdown correct, deleting from hub")

	delete(h.Rooms, room.Id)
	log.Printf("room has been removed from hub")
}

// Redirect message to room broadcasting channel by id
func (h *Hub) RedirectMsg(msg Message) {
	h.Rooms[msg.Room].Broadcast <- msg
	log.Println("Msg redirected")
}

// Check if room exist in hub
func (h *Hub) CheckRoomInHub(roomId string) bool {
	if _, ok := h.Rooms[roomId]; ok {
		return true
	}

	return false
}

// Check if requested room in hub by id
func (h *Hub) CheckRoomRequest(roomId string) *Room {
	if !h.CheckRoomInHub(roomId) {
		log.Printf("%s room not in hub", roomId)

		return nil
	}
	log.Printf("%s room in hub", roomId)

	return h.Rooms[roomId]
}

func (h *Hub) GetRoomById(roomId string) *Room {
	room := h.CheckRoomRequest(roomId)
	if room != nil {
		return room
	}

	log.Printf("%s not in hub", roomId)
	return nil
}
