package lobby

import "github.com/hetonei/arcanery-go-backend/internal/domain/models"

type Player struct {
	user  models.User
	class models.Class
	hand  []models.Card
	deck  []models.Card
}

type Round struct {
	count uint
	timer int
	turn  Player
}

type Game struct {
	rounds    []Round
	playerOne Player
	playerTwo Player
}

type RoomConnection interface {
	Connect()
	Disconnect()
}

type Room struct {
	game       Game
	connection RoomConnection
}
