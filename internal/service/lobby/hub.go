package lobby

import "log"

type Message struct {
	Room  string      `json:"rid"`
	Event string      `json:"type"`
	Data  interface{} `json:"data"`
}

type Hub struct {
	Rooms      map[string]*Room
	Broadcast  chan Message
	Register   chan *Room
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

func CheckRoomRequest(roomId string) *Room {
	if !H.CheckRoomInHub(roomId) {
		log.Printf("%s room not in hub", roomId)
		return nil
	}
	log.Printf("%s room in hub", roomId)
	return H.Rooms[roomId]
}
