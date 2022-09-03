package lobby

func GetAllSubscribersID(roomId string) []string {
	var subsId []string
	room := GetRoomById(roomId)
	for _, sub := range room.Subs {
		subsId = append(subsId, sub.SubId)
	}

	return subsId
}

func EventHandler(m Message) Message {
	switch m.Event {
	case "connected":
		m.Data = GetAllSubscribersID(m.Room)
	case "chat":
	}

	return m
}
