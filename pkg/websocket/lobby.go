package websocket

type Lobby struct {
	// Registered clients
	Clients map[*Client]bool

	// Inbound messages
	Broadcast chan string

	// Register requests
	Register chan *Client

	// Unregister requests
	Unregister chan *Client

	Content string
}

func (l *Lobby) Run() {
	for {
		select {
		case c := <-l.Register:
			l.Clients[c] = true
			c.send <- []byte(l.Content)
			break

		case c := <-l.Unregister:
			_, ok := l.Clients[c]
			if ok {
				delete(l.Clients, c)
				close(c.send)
			}
			break

		case m := <-l.Broadcast:
			l.Content = m
			l.broadcastMessage()
			break
		}
	}
}

func (l *Lobby) broadcastMessage() {
	for c := range l.Clients {
		select {
		case c.send <- []byte(l.Content):
			break

		// We can't reach the client
		default:
			close(c.send)
			delete(l.Clients, c)
		}
	}
}
