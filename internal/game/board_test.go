package game

import (
	"testing"

	"kings-corner/internal/deck"

	"github.com/stretchr/testify/assert"
)

func setupBoard() *Board {
	d := deck.New()
	return New(d)
}

func TestNew(t *testing.T) {
	d := deck.New()
	b := New(d)

	assert.Equal(t, b.deck, d)
	assert.Equal(t, len(b.Field), int(FIELDS_NUMBER))
	assert.Empty(t, b.Players)
	assert.Equal(t, int(b.CurrentTurn), 0)
}

func TestJoin(t *testing.T) {
	b := setupBoard()

	player := NewPlayer()

	b.Join(player)

	assert.NotEmpty(t, b.Players)
	assert.Equal(t, b.Players[0].ID(), player.ID())
}

func TestDrawPlayersHand(t *testing.T) {
	b := setupBoard()

	b.Join(NewPlayer())
	b.Join(NewPlayer())
	b.Join(NewPlayer())

	b.drawPlayersHand()

	for i, p := range b.Players {
		kcPlayer := p.(*kcPlayer)

		if i == int(b.CurrentTurn) {
			assert.Equal(t, len(kcPlayer.hand), INITIAL_HAND_CARDS+1)
			continue
		}
		assert.Equal(t, len(kcPlayer.hand), INITIAL_HAND_CARDS)
	}
}

func TestBuildField(t *testing.T) {
	b := setupBoard()

	b.buildField()

	for i, f := range b.Field {
		if i < 4 {
			assert.NotEmpty(t, f)
		} else {
			assert.Empty(t, f)
		}
	}
}

func TestIsCorner(t *testing.T) {
	b := setupBoard()

	cornerValue := []uint8{4, 5, 6, 7}
	nonCornerValue := []uint8{0, 1, 2, 3}

	for _, c := range cornerValue {
		assert.True(t, b.isCorner(c))
	}

	for _, c := range nonCornerValue {
		assert.False(t, b.isCorner(c))
	}
}

func TestCheckFieldLevel(t *testing.T) {
	b := setupBoard()

	upperFieldLevel := uint8(FIELDS_NUMBER)

	for i := 0; i < FIELDS_NUMBER; i++ {
		assert.Nil(t, b.checkFieldLevel(uint8(i)))
	}

	assert.Error(t, b.checkFieldLevel(upperFieldLevel))
}

func TestSetNextTurn(t *testing.T) {
	b := setupBoard()

	b.Join(NewPlayer())
	b.Join(NewPlayer())

	b.setNextTurn()
	assert.Equal(t, b.CurrentTurn, uint8(1))

	b.setNextTurn()
	assert.Equal(t, b.CurrentTurn, uint8(0))

	b.setNextTurn()
	assert.Equal(t, b.CurrentTurn, uint8(1))
}

func TestDrawPlayerTurn(t *testing.T) {
	b := setupBoard()

	b.Join(NewPlayer())
	b.Join(NewPlayer())

	b.drawPlayerTurn()

	kcPlayer := b.Players[b.CurrentTurn].(*kcPlayer)
	assert.Len(t, kcPlayer.hand, 1)
}
