package game

import (
	"fmt"

	"kings-corner/deck"
)

const (
	INITIAL_HAND_CARDS = 7
	FIELDS_NUMBER      = 8
)

type Board struct {
	Deck     deck.Deck
	Field    [FIELDS_NUMBER][]deck.Card
	PlayTurn chan Turn

	players     []Player
	currentTurn uint8
}

func New(d deck.Deck) Board {
	return Board{d, [8][]deck.Card{}, make(chan Turn), []Player{}, 0}
}

func (b *Board) Join(p Player) {
	p.SetPlayTurn(b.PlayTurn)

	b.players = append(b.players, p)
}

func (b *Board) Init() {
	// maximum players treatment
	b.Deck.Shuffle()

	b.drawPlayersHand()

	b.buildField()

	b.run()
}

func (b *Board) drawPlayersHand() {
	totalCards := INITIAL_HAND_CARDS * len(b.players)

	for i := 0; i < totalCards; i++ {
		card := *b.Deck.Pop()
		playerSelection := i % len(b.players)

		b.players[playerSelection].Draw(card)
	}
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
		return NewFieldAccessError(FieldLevelDoesNotExist)
	}

	return nil
}

func (b *Board) run() {
	for {
		select {
		case t := <-b.PlayTurn:
			b.play(t)
		default:
			fmt.Println("End Game")
			return
		}
	}
}

func (b *Board) play(t Turn) error {
	t.Play(b)
	b.setNextTurn()

	return nil
}

func (b *Board) setNextTurn() {
	b.currentTurn++

	if int(b.currentTurn) == len(b.players) {
		b.currentTurn = 0
	}
}
