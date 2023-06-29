package game

import (
	"testing"

	"kings-corner/internal/deck"

	"github.com/stretchr/testify/assert"
)

func setupBoardWithPlayers(p Player) *Board {
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

	firstTurn := turn.validateTurn(*b)
	assert.Error(t, firstTurn)

	b.setNextTurn()

	secondTurn := turn.validateTurn(*b)
	assert.Error(t, secondTurn)

	b.setNextTurn()
	playerTurn := turn.validateTurn(*b)
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
		newPlayedCardError(p.ID(), DifferentColorCard).Error(),
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
		newPlayedCardError(p.ID(), OneRankLower).Error(),
	)
}

func TestCardTurnPlay(t *testing.T) {
	p := NewPlayer()
	b := setupBoardWithPlayers(p)

	fieldLevel := 0
	fieldCards := b.Field[fieldLevel]
	validCard := getNextValidCard(fieldCards[0])

	b.Players[2].draw(validCard)

	ct := &CardTurn{
		FieldLevel: uint8(fieldLevel),
		Card:       validCard,
		turn:       turn{p},
	}

	err := ct.Play(b)
	assert.Error(t, err, "isn't player turn")
	assert.IsType(t, err, &playerAccessError{})

	b.setNextTurn()
	b.setNextTurn()
	ct.FieldLevel = 9
	err = ct.Play(b)
	assert.Error(t, err, "invalid field access")
	assert.IsType(t, err, &fieldAccessError{})

	fieldLevelFulfilled := uint8(fieldLevel) + 1
	b.Field[fieldLevelFulfilled] = make([]deck.Card, 13)
	ct.FieldLevel = fieldLevelFulfilled
	err = ct.Play(b)
	assert.Error(t, err, "field level fulfilled")
	assert.IsType(t, err, &fieldAccessError{})

	ct.FieldLevel = uint8(fieldLevel)
	ct.Card = fieldCards[0]
	err = ct.Play(b)
	assert.Error(t, err, "bad card played")
	assert.IsType(t, err, &playedCardError{})

	fieldLevelEmpty := uint8(fieldLevel) + 1
	b.Field[fieldLevelEmpty] = []deck.Card{}
	ct.FieldLevel = fieldLevelEmpty
	ct.Card = deck.Card{Suit: deck.Diamond, Rank: deck.King}
	err = ct.Play(b)
	assert.Error(t, err, "king on non corners")
	assert.IsType(t, err, &playedCardError{})

	ct.Card = validCard
	ct.FieldLevel = uint8(fieldLevel)
	err = ct.Play(b)
	assert.Nil(t, err, "played normal card")
	assertPlayedCardTurn(t, p, validCard, b.Field[ct.FieldLevel])

	kingCard := deck.Card{Suit: deck.Diamond, Rank: deck.King}
	p.draw(kingCard)
	ct.Card = kingCard
	ct.FieldLevel = 4
	err = ct.Play(b)
	assert.Nil(t, err, "played king card on corner")
	assertPlayedCardTurn(t, p, kingCard, b.Field[ct.FieldLevel])
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

func TestMoveTurnPlay(t *testing.T) {
	p := NewPlayer()
	b := setupBoardWithPlayers(p)

	fieldLevel := [2]uint8{0, 0}
	moveToFieldLevel := 1

	b.Field[0][0] = deck.Card{Suit: deck.Club, Rank: deck.Two}
	b.Field[0] = append(b.Field[0], deck.Card{Suit: deck.Diamond, Rank: deck.Ace})
	b.Field[moveToFieldLevel][0] = deck.Card{Suit: deck.Heart, Rank: deck.Three}

	mt := &MoveTurn{
		FieldLevel:       fieldLevel,
		MoveToFieldLevel: uint8(moveToFieldLevel),
		turn:             turn{p},
	}

	err := mt.Play(b)
	assert.Error(t, err, "isn't player turn")
	assert.IsType(t, err, &playerAccessError{})

	b.setNextTurn()
	b.setNextTurn()
	err = mt.Play(b)
	assert.Nil(t, err)
	assert.Len(t, b.Field[0], 0)
	assert.Len(t, b.Field[moveToFieldLevel], 3)

	invalidFieldLevel := [2]uint8{9, 1}
	mt.FieldLevel = invalidFieldLevel
	err = mt.Play(b)
	assert.Error(t, err, "invalid field access")
	assert.IsType(t, err, &fieldAccessError{})

	invalidMoveToFieldLevel := 9
	mt.MoveToFieldLevel = uint8(invalidMoveToFieldLevel)
	err = mt.Play(b)
	assert.Error(t, err, "invalid move to field access")
	assert.IsType(t, err, &fieldAccessError{})

	mt.FieldLevel = fieldLevel
	mt.MoveToFieldLevel = uint8(moveToFieldLevel)
	err = mt.Play(b)
	assert.Error(t, err, "no card on selected field")
	assert.IsType(t, err, &fieldAccessError{})

	mt.FieldLevel = [2]uint8{0, 1}
	mt.MoveToFieldLevel = uint8(moveToFieldLevel)
	err = mt.Play(b)
	assert.Error(t, err, "invalid card field index")
	assert.IsType(t, err, &fieldAccessError{})

	b.Field[0] = append(b.Field[0], deck.Card{Suit: deck.Club, Rank: deck.Four})
	mt.FieldLevel = fieldLevel
	mt.MoveToFieldLevel = uint8(moveToFieldLevel)
	err = mt.Play(b)
	assert.Error(t, err, "invalid card validation")
	assert.IsType(t, err, &playedCardError{})
}

func TestPassTurnPlay(t *testing.T) {
	p := NewPlayer()
	b := setupBoardWithPlayers(p)

	pt := &PassTurn{
		turn: turn{p},
	}

	err := pt.Play(b)
	assert.Error(t, err, "isn't player turn")
	assert.IsType(t, err, &playerAccessError{})

	b.setNextTurn()
	b.setNextTurn()

	err = pt.Play(b)
	assert.Nil(t, err)
}
