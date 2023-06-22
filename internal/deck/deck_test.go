package deck_test

import (
	"fmt"
	"testing"

	"kings-corner/internal/deck"

	"github.com/stretchr/testify/assert"
)

const LEN_DECK = 52

func TestNewDeck(t *testing.T) {
	d := deck.New()

	assert.Equal(t, LEN_DECK, len(d), fmt.Sprintf("should have %d cards on deck", LEN_DECK))
}

func TestShuffle(t *testing.T) {
	d := deck.New()

	shuffleDeck := deck.New()
	shuffleDeck.Shuffle()

	different := 0
	equal := 0

	for i := range shuffleDeck {
		if shuffleDeck[i] == d[i] {
			equal++
			continue
		}
		different++
	}

	assert.Greater(t, different, equal)
}

func TestPop(t *testing.T) {
	d := deck.New()

	poppedCard := d.Pop()

	assert.Equal(t, len(d), LEN_DECK-1)

	for i := range d {
		assert.NotEqual(t, d[i], poppedCard)
	}
}

func TestPopLastCard(t *testing.T) {
	d := deck.New()

	for range d {
		card := d.Pop()
		assert.NotNil(t, card)
	}

	card := d.Pop()
	assert.Nil(t, card)
}

func TestPopNoKing(t *testing.T) {
	d := deck.New()

	d[12], d[0] = d[0], d[12]

	card := d.PopNoKing()

	assert.NotEqual(t, card.Rank, deck.King)
}

func TestSuitColor(t *testing.T) {
	suits := [...]deck.Suit{deck.Spade, deck.Diamond, deck.Club, deck.Heart}

	for _, s := range suits {
		color := s.Color()

		if s == deck.Diamond || s == deck.Heart {
			assert.Equal(t, color, deck.Red)
		} else {
			assert.Equal(t, color, deck.Black)
		}
	}
}

func TestCardIsEqual(t *testing.T) {
	card := deck.Card{
		Suit: deck.Heart,
		Rank: deck.King,
	}

	differentCard := deck.Card{
		Suit: deck.Heart,
		Rank: deck.Jack,
	}

	equalCard := deck.Card{
		Suit: deck.Heart,
		Rank: deck.King,
	}

	assert.False(t, card.IsEqual(differentCard))
	assert.True(t, card.IsEqual(equalCard))
}

func TestCardIsOneRankHigherThan(t *testing.T) {
	card := deck.Card{
		Suit: deck.Heart,
		Rank: deck.Queen,
	}

	higherRank := deck.Card{
		Suit: deck.Heart,
		Rank: deck.King,
	}

	lowerRank := deck.Card{
		Suit: deck.Heart,
		Rank: deck.Jack,
	}

	assert.False(t, card.IsOneRankHigherThan(higherRank))
	assert.True(t, card.IsOneRankHigherThan(lowerRank))
}

func TestCardIsSameColor(t *testing.T) {
	card := deck.Card{
		Suit: deck.Heart,
		Rank: deck.Queen,
	}

	sameColor := deck.Card{
		Suit: deck.Diamond,
		Rank: deck.Queen,
	}

	notSameColor := deck.Card{
		Suit: deck.Club,
		Rank: deck.Queen,
	}

	assert.False(t, card.IsSameColor(notSameColor))
	assert.True(t, card.IsSameColor(sameColor))
}
