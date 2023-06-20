package game

import (
	"kings-corner/deck"

	"github.com/rs/xid"
)

func NewPlayer() Player {
	return &kcPlayer{
		id:       xid.New().String(),
		playTurn: nil,
	}
}

type Player interface {
	ID() string
	Draw(card deck.Card)
	Play(card deck.Card)
	SetPlayTurn(chan<- Turn)
	IsWinner() bool
}

type kcPlayer struct {
	id   string
	hand []deck.Card

	playTurn chan<- Turn
}

func (p *kcPlayer) ID() string {
	return p.id
}

func (p *kcPlayer) Draw(card deck.Card) {
	p.hand = append(p.hand, card)
}

func (p *kcPlayer) Play(card deck.Card) {
	newHand := []deck.Card{}
	for i := range p.hand {
		if p.hand[i].IsEqual(card) {
			continue
		}

		newHand = append(newHand, p.hand[i])
	}

	p.hand = newHand
}

func (p *kcPlayer) SetPlayTurn(playTurn chan<- Turn) {
	p.playTurn = playTurn
}

func (p *kcPlayer) IsWinner() bool {
	return len(p.hand) == 0
}
