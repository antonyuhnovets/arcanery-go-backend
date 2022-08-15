package websocket

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

var upgrader = &ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ReadMsgsWS(w http.ResponseWriter, r *http.Request, h http.Header) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Can`t connect (%s)", err)
		return
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read msg")
			break
		}
		conn.WriteMessage(t, msg)
	}
}
