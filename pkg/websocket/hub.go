package websocket

type Message struct {
	Data []byte
	Room string
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
			h.AddRoom(r)
		case r := <-h.Unregister:
			h.RemoveRoom(r)
		case m := <-h.Broadcast:
			h.RedirectMsg(m)
		}
	}
}

func (h *Hub) AddRoom(room *Room) {
	h.Rooms[room.Id] = room
	go room.Start()
}

func (h *Hub) RemoveRoom(room *Room) {
	room.Shutdown()
	delete(h.Rooms, room.Id)
}

func (h *Hub) RedirectMsg(msg Message) {
	h.Rooms[msg.Room].ProcessMsg(msg)
}

func (h *Hub) CheckRoomInHub(roomId string) bool {
	if _, ok := h.Rooms[roomId]; ok {
		return true
	}
	return false
}
