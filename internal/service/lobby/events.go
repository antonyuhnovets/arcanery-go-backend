// Package represent lobby realisation. Main logic of hub with rooms, message and events processing.
// File implement event handler.

package lobby

import (
	"encoding/json"
	"log"
)

// Handle event from msg.
func (h *Hub) HandleMsg(m Message) Message {
	switch m.Event {
	case "connected": // new subscriber in room event
		id := m.PullData("uid")
		m.Data = id

	case "chat": // text message
	}

	return m
}

// Pull some information from message
func (m Message) PullData(data string) interface{} {
	b, err := json.Marshal(m.Data)
	if err != nil {
		log.Println(err)
	}

	value := GetDataObj(b, data)

	return value
}

// Find value of needed field
func GetDataObj(data []byte, key string) interface{} {
	d := make(map[string]interface{})
	err := json.Unmarshal(data, &d)
	if err != nil {
		log.Println(err)
	}

	return d[key]
}

// Get list of all subscribers ids in room
func GetAllSubscribersID(roomId string) []string {
	var subsId []string
	room := GetRoomById(roomId)
	for _, sub := range room.Subs {
		subsId = append(subsId, sub.Id)
	}

	return subsId
}
