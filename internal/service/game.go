package service

import "github.com/hetonei/arcanery-go-backend/internal/domain/models"

type Player struct {
	user  models.User
	class models.Class
	hand  []models.Card
	deck  []models.Card
}

type Round struct {
	count int
	timer int
	turn  Player
}

type Game struct {
	rounds    []Round
	playerOne Player
	playerTwo Player
}

type GameService interface {
	StartGame()
}

func (g Game) StartGame() {
	return
}

func GetGameService(p1, p2 Player) GameService {
	rounds := make([]Round, 10)
	for c := range rounds {
		var turn Player
		if c%2 == 0 {
			turn = p2
		} else {
			turn = p1
		}
		rounds[c] = Round{
			count: c,
			timer: 30,
			turn:  turn,
		}
	}
	return Game{
		rounds:    rounds,
		playerOne: p1,
		playerTwo: p2,
	}
}
