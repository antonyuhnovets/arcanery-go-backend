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
	subId string
	Room  string
	Conn  *Connection
}

func ServeWs(w http.ResponseWriter, r *http.Request, roomId string) {
	room := CheckRoomRequest(roomId)

	ws, err := UpgradeToWs(w, r)
	if err != nil {
		log.Fatal(err)
	}

	c := &Connection{send: make(chan []byte, 256), ws: ws}
	s := Subscription{r.RemoteAddr, roomId, c}

	s.Subscribe(room)
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
	var room *Room

	if !H.CheckRoomInHub(roomId) {
		room = CreateRoom(roomId)
		go room.Start()
		H.Register <- room
	} else {
		room = H.Rooms[roomId]
	}

	return room
}

func (s Subscription) Subscribe(room *Room) {
	room.Register <- s

	s.SetOpts()

	go s.writePump()
	go s.readPump(room)
}

func (s Subscription) readPump(room *Room) {
	c := s.Conn

	defer func() {
		room.Unregister <- s
		c.ws.Close()
	}()

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := Message{msg, s.Room}
		H.Broadcast <- m
	}
}

func (s *Subscription) writePump() {
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
	return c.ws.WriteMessage(mt, payload)
}

func (s Subscription) SetOpts() {
	c := s.Conn

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
}
