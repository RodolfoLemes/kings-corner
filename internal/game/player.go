package game

import (
	"kings-corner/internal/deck"

	"github.com/rs/xid"
)

type Play func(Turn) error

func NewPlayer() Player {
	return &kcPlayer{
		id: xid.New().String(),
	}
}

type Player interface {
	ID() string
	Play(Turn) error
	Hand() []deck.Card
	draw(card deck.Card)
	withdraw(card deck.Card)
	setPlay(Play)
	isWinner() bool
}

type kcPlayer struct {
	id   string
	hand []deck.Card

	playFunc Play
}

func (p *kcPlayer) ID() string {
	return p.id
}

func (p *kcPlayer) Hand() []deck.Card {
	return p.hand
}

func (p *kcPlayer) draw(card deck.Card) {
	p.hand = append(p.hand, card)
}

func (p *kcPlayer) withdraw(card deck.Card) {
	newHand := []deck.Card{}
	for i := range p.hand {
		if p.hand[i].IsEqual(card) {
			continue
		}

		newHand = append(newHand, p.hand[i])
	}

	p.hand = newHand
}

func (p *kcPlayer) Play(t Turn) error {
	t.setPlayer(p)
	return p.playFunc(t)
}

func (p *kcPlayer) setPlay(playFunc Play) {
	p.playFunc = playFunc
}

func (p *kcPlayer) isWinner() bool {
	return len(p.hand) == 0
}
