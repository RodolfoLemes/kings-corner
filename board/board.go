package board

import (
	"fmt"

	deck "kings-corner"
)

const INITIAL_HAND_CARDS = 7

type PlayTurn struct {
	FieldLevel uint8
	card       deck.Card
	player     Player
}

type Board struct {
	Deck     deck.Deck
	Field    [8][]deck.Card
	PlayTurn chan *PlayTurn

	players     []Player
	currentTurn uint8
}

func New(d deck.Deck) Board {
	return Board{d, [8][]deck.Card{}, make(chan *PlayTurn), []Player{}, 0}
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

		card := *b.Deck.Pop()

		b.Field[fieldSelection] = []deck.Card{card}
	}
}

func (b *Board) run() {
	for {
		select {
		case pt := <-b.PlayTurn:
			b.play(pt)
		default:
			fmt.Println("End Game")
			return
		}
	}
}

func (b *Board) play(pt *PlayTurn) error {
	if pt.player.ID() != b.players[b.currentTurn].ID() {
		return fmt.Errorf("invalid player access")
	}

	if pt.FieldLevel > 8 {
		return fmt.Errorf("invalid field access")
	}

	selectedField := &b.Field[pt.FieldLevel]

	if len(*selectedField) > int(deck.MaxRank) {
		return fmt.Errorf("field level already fulfilled")
	}

	lastFieldCard := (*selectedField)[len(*selectedField)-1]

	if lastFieldCard.IsSameColor(pt.card) {
		return fmt.Errorf("invalid played card, must be from different color")
	}

	if !lastFieldCard.IsOneRankHigherThan(pt.card) {
		return fmt.Errorf("invalid played card, must be one rank lower")
	}

	*selectedField = append(*selectedField, pt.card)

	pt.player.Play(pt.card)

	b.setNextTurn()

	return nil
}

func (b *Board) setNextTurn() {
	b.currentTurn++

	if int(b.currentTurn) == len(b.players) {
		b.currentTurn = 0
	}
}
