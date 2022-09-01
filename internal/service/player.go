package service

import (
	"github.com/hetonei/arcanery-go-backend/internal/domain/models"
)

type Player struct {
	id        string
	ready     bool
	cardsPool map[string]models.Card
	hand      map[string]models.Card
}

func (p *Player) GetReady() {
	p.ready = true
}

func (p *Player) ChoseCards(cards ...string) {
	for _, card := range cards {
		p.hand[card] = p.cardsPool[card]
	}
}

func (p *Player) PlayCard(id string) {
}

func (p *Player) ChangeCard(oldCard, newCard string) {
}
