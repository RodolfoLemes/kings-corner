//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Color uint8

const (
	Red Color = iota
	Black
)

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	// Joker // this is a special case
)

func (s Suit) Color() Color {
	if s == Diamond || s == Heart {
		return Red
	}

	return Black
}

var suits = [...]Suit{Spade, Diamond, Club, Heart}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	MinRank = Ace
	MaxRank = King
)

type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	/* if c.Suit == Joker {
		return c.Suit.String()
	} */
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

func (c Card) IsEqual(cc Card) bool {
	return c.Rank == cc.Rank && c.Suit == cc.Suit
}

func (c Card) IsSameColor(cc Card) bool {
	return c.Color() == cc.Color()
}

func (c Card) IsOneRankHigherThan(cc Card) bool {
	return c.Rank-1 == cc.Rank
}

type Deck []Card

func New() Deck {
	var cards Deck
	for _, suit := range suits {
		for rank := MinRank; rank <= MaxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	return cards
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(*d), func(i, j int) {
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	})
}

func (d *Deck) Pop() *Card {
	if len(*d) == 0 {
		return nil
	}

	card := (*d)[0]

	*d = (*d)[1:]

	return &card
}

func (d *Deck) PopNoKing() *Card {
	card := d.Pop()

	if card.Rank == King {
		*d = append(*d, *card)
		card := d.PopNoKing()
		return card
	}

	return card
}
