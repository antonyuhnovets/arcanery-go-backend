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
			m := h.HandleMsg(msg)
			h.RedirectMsg(m)
		}
	}
}

// Send room registration request to hub
func RegisterRoom(r *Room) {
	H.Register <- r
}

// Send room unregister request to hub
func Unregister(r *Room) {
	H.Unregister <- r
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
		log.Println(true)
		return true
	}
	return false
}

// Check if requested room in hub by id
func CheckRoomRequest(roomId string) *Room {
	if !H.CheckRoomInHub(roomId) {
		log.Printf("%s room not in hub", roomId)
		return nil
	}
	log.Printf("%s room in hub", roomId)
	return H.Rooms[roomId]
}
