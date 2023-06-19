package game

import (
	"testing"

	"kings-corner/deck"

	"github.com/stretchr/testify/assert"
)

func setupBoardWithPlayers(p Player) Board {
	b := setupBoard()

	b.Join(NewPlayer())
	b.Join(NewPlayer())
	b.Join(p)

	b.drawPlayersHand()
	b.buildField()

	return b
}

func TestValidateTurn(t *testing.T) {
	p := NewPlayer()

	b := setupBoardWithPlayers(p)

	turn := turn{p}

	firstTurn := turn.validateTurn(b)
	assert.Error(t, firstTurn)

	b.setNextTurn()

	secondTurn := turn.validateTurn(b)
	assert.Error(t, secondTurn)

	b.setNextTurn()
	playerTurn := turn.validateTurn(b)
	assert.Nil(t, playerTurn)
}

func TestValidateCardInsertion(t *testing.T) {
	p := NewPlayer()
	turn := turn{p}

	initialCard := deck.Card{
		Suit: deck.Diamond,
		Rank: deck.Queen,
	}
	insertionCard := deck.Card{
		Suit: deck.Club,
		Rank: deck.Jack,
	}

	assert.Nil(t, turn.validateCardInsertion(initialCard, insertionCard))
}

func TestValidateCardInsertionDifferentColor(t *testing.T) {
	p := NewPlayer()
	turn := turn{p}

	initialCard := deck.Card{
		Suit: deck.Diamond,
		Rank: deck.Queen,
	}
	insertionCard := deck.Card{
		Suit: deck.Diamond,
		Rank: deck.Queen,
	}

	differentColorErr := turn.validateCardInsertion(initialCard, insertionCard)

	assert.Error(t, differentColorErr)
	assert.EqualError(
		t,
		differentColorErr,
		NewPlayedCardError(p.ID(), DifferentColorCard).Error(),
	)
}

func TestValidateCardInsertionNotOneRankHigherThan(t *testing.T) {
	p := NewPlayer()
	turn := turn{p}

	initialCard := deck.Card{
		Suit: deck.Diamond,
		Rank: deck.Queen,
	}
	insertionCard := deck.Card{
		Suit: deck.Club,
		Rank: deck.King,
	}

	notOneRankHigher := turn.validateCardInsertion(initialCard, insertionCard)

	assert.Error(t, notOneRankHigher)
	assert.EqualError(
		t,
		notOneRankHigher,
		NewPlayedCardError(p.ID(), OneRankLower).Error(),
	)
}

func TestCardTurnPlay(t *testing.T) {
	p := NewPlayer()
	b := setupBoardWithPlayers(p)

	fieldLevel := 0
	fieldCards := b.Field[fieldLevel]
	validCard := getNextValidCard(fieldCards[0])

	b.players[2].Draw(validCard)

	ct := &cardTurn{
		fieldLevel: uint8(fieldLevel),
		card:       validCard,
		turn:       turn{p},
	}

	err := ct.Play(&b)
	assert.Error(t, err, "isn't player turn")
	assert.IsType(t, err, &PlayerAccessError{})

	b.setNextTurn()
	b.setNextTurn()
	ct.fieldLevel = 9
	err = ct.Play(&b)
	assert.Error(t, err, "invalid field access")
	assert.IsType(t, err, &FieldAccessError{})

	fieldLevelFulfilled := uint8(fieldLevel) + 1
	b.Field[fieldLevelFulfilled] = make([]deck.Card, 13)
	ct.fieldLevel = fieldLevelFulfilled
	err = ct.Play(&b)
	assert.Error(t, err, "field level fulfilled")
	assert.IsType(t, err, &FieldAccessError{})

	ct.fieldLevel = uint8(fieldLevel)
	ct.card = fieldCards[0]
	err = ct.Play(&b)
	assert.Error(t, err, "bad card played")
	assert.IsType(t, err, &PlayedCardError{})

	fieldLevelEmpty := uint8(fieldLevel) + 1
	b.Field[fieldLevelEmpty] = []deck.Card{}
	ct.fieldLevel = fieldLevelEmpty
	ct.card = deck.Card{Suit: deck.Diamond, Rank: deck.King}
	err = ct.Play(&b)
	assert.Error(t, err, "king on non corners")
	assert.IsType(t, err, &PlayedCardError{})

	ct.card = validCard
	ct.fieldLevel = uint8(fieldLevel)
	err = ct.Play(&b)
	assert.Nil(t, err, "played normal card")
	assertPlayedCardTurn(t, p, validCard, b.Field[ct.fieldLevel])

	kingCard := deck.Card{Suit: deck.Diamond, Rank: deck.King}
	p.Draw(kingCard)
	ct.card = kingCard
	ct.fieldLevel = 4
	err = ct.Play(&b)
	assert.Nil(t, err, "played king card on corner")
	assertPlayedCardTurn(t, p, kingCard, b.Field[ct.fieldLevel])
}

func getNextValidCard(c deck.Card) deck.Card {
	suit := c.Suit + 1
	if suit > 3 {
		suit = 0
	}

	return deck.Card{
		Suit: suit,
		Rank: c.Rank - 1,
	}
}

func assertPlayedCardTurn(t *testing.T, p Player, validCard deck.Card, fieldCards []deck.Card) {
	playedCardFieldExists := false
	playedCardHandExists := false

	for _, f := range fieldCards {
		if f.IsEqual(validCard) {
			playedCardFieldExists = true
		}
	}
	assert.True(t, playedCardFieldExists)

	kcPlayer := p.(*kcPlayer)
	for _, h := range kcPlayer.hand {
		if h.IsEqual(validCard) {
			playedCardHandExists = true
		}
	}
	assert.False(t, playedCardHandExists)
}