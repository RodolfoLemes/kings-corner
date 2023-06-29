package game

import (
	"math/rand"
	"testing"
	"time"

	"kings-corner/internal/deck"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewPLayer(t *testing.T) {
	player := NewPlayer()

	assert.IsType(t, player, &kcPlayer{})

	kcPLayer := player.(*kcPlayer)

	assert.Empty(t, kcPLayer.hand)
}

func TestDraw(t *testing.T) {
	player := NewPlayer()

	card := generateRandomCard()

	player.draw(card)

	kcPLayer := player.(*kcPlayer)

	assert.True(t, card.IsEqual(kcPLayer.hand[0]))
}

func TestWithdraw(t *testing.T) {
	player := NewPlayer()

	card := generateRandomCard()

	player.draw(card)

	player.withdraw(card)

	kcPLayer := player.(*kcPlayer)

	assert.Equal(t, len(kcPLayer.hand), 0)
}

func generateRandomCard() deck.Card {
	rankCard1 := rand.Intn(int(deck.MaxRank) - 1)
	suitCard1 := rand.Intn(int(deck.Heart))

	return deck.Card{
		Suit: deck.Suit(suitCard1),
		Rank: deck.Rank(rankCard1),
	}
}
