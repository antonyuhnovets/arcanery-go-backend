// Package manage websocket connection, handshake using gorilla lib.

package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Options for websocket connection.
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type Connection struct {
	Id   string
	ws   *websocket.Conn
	send chan interface{} // channel storing outcoming messages
}

func MakeConnection(ws *websocket.Conn) *Connection {
	c := Connection{
		send: make(chan interface{}),
		ws:   ws,
	}
	return &c
}

// Upgrades connection protocol to websocket
func UpgradeToWs(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize,
		WriteBufferSize: maxMessageSize,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func (c *Connection) SendMsg(msg interface{}) {
	c.send <- msg
}

// Read message from websocket connection and handle it with passed function
func (c *Connection) ReadPump(msgHandler func([]byte)) {
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
		msgHandler(msg)
	}
}

// Check if there is new outcoming message
// Marshal it to bytes and pass to write function
// Pass an empty msg if no msgs sended
func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			bMsg, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}

			if err := c.write(websocket.TextMessage, bMsg); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Connection) Close() {
	close(c.send)
	c.ws.Close()
}

func (c *Connection) SetConnectionId(id string) {
	c.Id = id
}

func (c *Connection) GetConnectionId() string {
	return c.Id
}

// Writes a message with the given message type and payload
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))

	return c.ws.WriteMessage(mt, payload)
}
