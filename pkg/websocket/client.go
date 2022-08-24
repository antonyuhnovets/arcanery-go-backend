package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type Connection struct {
	ws   *websocket.Conn
	send chan []byte // Channel storing outcoming messages
}

type Subscription struct {
	SubId string
	Room  string
	Conn  *Connection
}

func ServeWs(w http.ResponseWriter, r *http.Request, roomId string) {
	room := CheckRoomRequest(roomId)

	ws, err := UpgradeToWs(w, r)
	if err != nil {
		log.Fatal(err)
	}

	s := CreateSubscription(r.RemoteAddr, roomId, ws)

	s.Subscribe(room)
}

func CreateSubscription(id, roomId string, ws *websocket.Conn) Subscription {
	c := &Connection{
		send: make(chan []byte, 256),
		ws:   ws,
	}

	return Subscription{id, roomId, c}
}

func UpgradeToWs(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func CheckRoomRequest(roomId string) *Room {
	if !H.CheckRoomInHub(roomId) {
		log.Printf("%s room not in hub", roomId)
		return nil
	}
	log.Printf("%s room in hub", roomId)
	return H.Rooms[roomId]
}

func (s Subscription) Subscribe(room *Room) {
	room.Register <- s

	go s.writePump(room)
	go s.readPump(room)
}

func (s Subscription) readPump(room *Room) {
	c := s.Conn

	defer func() {
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		H.Broadcast <- Message{msg, room.Id}
	}
}

func (s Subscription) writePump(room *Room) {
	c := s.Conn

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))

	return c.ws.WriteMessage(mt, payload)
}
