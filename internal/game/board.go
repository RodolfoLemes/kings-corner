package game

import (
	"fmt"

	"kings-corner/internal/deck"
	"kings-corner/internal/pubsub"

	"github.com/rs/xid"
)

const (
	INITIAL_HAND_CARDS = 7
	FIELDS_NUMBER      = 8
)

type Board struct {
	ID          string
	CurrentTurn uint8
	IsStarted   bool
	Field       [FIELDS_NUMBER][]deck.Card

	deck    deck.Deck
	Players []Player

	pubsub pubsub.PubSub[Board]
}

func New(d deck.Deck) *Board {
	return &Board{
		ID:      xid.New().String(),
		Field:   [FIELDS_NUMBER][]deck.Card{},
		deck:    d,
		Players: []Player{},
		pubsub:  pubsub.New[Board](),
	}
}

func (b *Board) Listen() Board {
	return <-b.pubsub.Subscribe(b.channel())
}

func (b Board) channel() string {
	return fmt.Sprintf("game:%s", b.ID)
}

func (b *Board) Join(p Player) {
	b.Players = append(b.Players, p)
}

func (b *Board) Run() error {
	// maximum Players treatment
	b.deck.Shuffle()

	b.drawPlayersHand()

	b.buildField()

	return b.run()
}

func (b *Board) drawPlayersHand() {
	totalCards := INITIAL_HAND_CARDS * len(b.Players)

	for i := 0; i < totalCards; i++ {
		card := *b.deck.Pop()
		PlayerSelection := i % len(b.Players)

		b.Players[PlayerSelection].draw(card)
	}

	b.drawPlayerTurn()
}

func (b *Board) buildField() {
	noCornerFieldsNumber := 4

	for i := 0; i < noCornerFieldsNumber; i++ {
		fieldSelection := i % noCornerFieldsNumber

		card := *b.deck.PopNoKing()

		b.Field[fieldSelection] = []deck.Card{card}
	}
}

func (b *Board) isCorner(fieldLevel uint8) bool {
	return fieldLevel > 3
}

func (b *Board) checkFieldLevel(fieldLevel uint8) error {
	if fieldLevel > FIELDS_NUMBER-1 {
		return newFieldAccessError(FieldLevelDoesNotExist)
	}

	return nil
}

func (b *Board) setNextTurn() {
	b.CurrentTurn++

	if int(b.CurrentTurn) == len(b.Players) {
		b.CurrentTurn = 0
	}
}

func (b *Board) drawPlayerTurn() {
	card := b.deck.Pop()
	b.Players[b.CurrentTurn].draw(*card)
}

func (b *Board) run() error {
	b.IsStarted = true
	return b.pubsub.Publish("game", *b)
}

func (b *Board) play(t Turn) error {
	err := t.Play(b)
	if err != nil {
		return err
	}

	if b.hasWinner() {
		b.endGame() // TODO
	}

	err = b.pubsub.Publish(b.channel(), *b)

	return err
}

func (b *Board) hasWinner() bool {
	for _, p := range b.Players {
		if p.isWinner() {
			return true
		}
	}

	return false
}

func (b *Board) endGame() {}
