package game

import (
	"kings-corner/internal/deck"
	"kings-corner/internal/pubsub"
)

const (
	INITIAL_HAND_CARDS = 7
	FIELDS_NUMBER      = 8
)

type Board struct {
	Deck     deck.Deck
	Field    [FIELDS_NUMBER][]deck.Card
	PlayTurn chan Turn

	Players     []Player
	CurrentTurn uint8

	pubsub pubsub.PubSub[Board]
}

func New(d deck.Deck) Board {
	return Board{d, [8][]deck.Card{}, make(chan Turn), []Player{}, 0, pubsub.New[Board]()}
}

func (b *Board) Listen() Board {
	return <-b.pubsub.Subscribe("game")
}

func (b *Board) Join(p Player) {
	p.setPlayTurn(b.PlayTurn)

	b.Players = append(b.Players, p)
}

func (b *Board) Run() {
	// maximum players treatment
	b.Deck.Shuffle()

	b.drawPlayersHand()

	b.buildField()

	b.run()
}

func (b *Board) drawPlayersHand() {
	totalCards := INITIAL_HAND_CARDS * len(b.Players)

	for i := 0; i < totalCards; i++ {
		card := *b.Deck.Pop()
		playerSelection := i % len(b.Players)

		b.Players[playerSelection].Draw(card)
	}

	b.drawPlayerTurn()
}

func (b *Board) buildField() {
	noCornerFieldsNumber := 4

	for i := 0; i < noCornerFieldsNumber; i++ {
		fieldSelection := i % noCornerFieldsNumber

		card := *b.Deck.PopNoKing()

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
	card := b.Deck.Pop()
	b.Players[b.CurrentTurn].Draw(*card)
}

func (b *Board) run() {
	b.pubsub.Publish("game", *b)

	for {
		select {
		case t := <-b.PlayTurn:
			b.play(t)
		}
	}
}

func (b *Board) play(t Turn) {
	err := t.Play(b)
	if err != nil {
		return
	}

	if b.hasWinner() {
		b.endGame()
	}

	b.pubsub.Publish("game", *b)

	return
}

func (b *Board) hasWinner() bool {
	for _, p := range b.Players {
		if p.IsWinner() {
			return true
		}
	}

	return false
}

func (b *Board) endGame() {
	close(b.PlayTurn)
}
