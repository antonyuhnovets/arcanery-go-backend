package lobby

import (
	"encoding/json"
	"log"
	// "github.com/hetonei/arcanery-go-backend/internal/service"
)

func (h *Hub) HandleMsg(m Message) {
	room := GetRoomById(m.Room)
	switch m.Event {
	case "connected":
		id := m.PullData("uid")
		m.Data = id

		for k, sub := range room.Subs {
			log.Println(k, " ", sub)
		}
	case "chat":
	}
}

func (m Message) PullData(data string) interface{} {
	b, err := json.Marshal(m.Data)
	if err != nil {
		log.Println(err)
	}
	id := GetDataObj(b, data)
	return id
}

func GetAllSubscribersID(roomId string) []string {
	var subsId []string
	room := GetRoomById(roomId)
	for _, sub := range room.Subs {
		subsId = append(subsId, sub.SubId)
	}

	return subsId
}

func GetDataObj(data []byte, key string) interface{} {
	d := make(map[string]interface{})
	err := json.Unmarshal(data, &d)
	if err != nil {
		log.Println(err)
	}

	return d[key]
}
