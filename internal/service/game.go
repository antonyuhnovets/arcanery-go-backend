package service

import (
	// "encoding/json"

	"github.com/hetonei/arcanery-go-backend/internal/domain/models"
	// "go.mongodb.org/mongo-driver/event"
)

type Action struct {
	actionType string
	handler    func()
}

type Deck struct {
	Player1 []models.Card
	Player2 []models.Card
}

type Stage struct {
	actionPlayer1 chan Action
	actionPlayer2 chan Action
	done          chan bool
}

type Round struct {
	done chan bool
}

type Game struct {
	Player1 *Player
	Player2 *Player
}
