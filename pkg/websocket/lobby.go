package websocket

type Message struct {
	Data []byte
	Room string
}

type Subscription struct {
	Conn *Connection
	Room string
}

type Lobby struct {
	// Registered connections.
	Rooms map[string]map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan Message

	// Register requests from the connections.
	Register chan Subscription

	// Unregister requests from connections.
	Unregister chan Subscription
}

func NewLobby() *Lobby {
	l := &Lobby{
		Broadcast:  make(chan Message),
		Register:   make(chan Subscription),
		Unregister: make(chan Subscription),
		Rooms:      make(map[string]map[*Connection]bool),
	}
	return l
}

func (l *Lobby) Run() {
	for {
		select {
		case s := <-l.Register:
			connections := l.Rooms[s.Room]
			if connections == nil {
				connections = make(map[*Connection]bool)
				l.Rooms[s.Room] = connections
			}
			l.Rooms[s.Room][s.Conn] = true
		case s := <-l.Unregister:
			connections := l.Rooms[s.Room]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.send)
					if len(connections) == 0 {
						delete(l.Rooms, s.Room)
					}
				}
			}
		case m := <-l.Broadcast:
			connections := l.Rooms[m.Room]
			for c := range connections {
				select {
				case c.send <- m.Data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(l.Rooms, m.Room)
					}
				}
			}
		}
	}
}
