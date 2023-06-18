//go:generate stringer -type=Suit,Rank

package deck

import "fmt"

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
	minRank = Ace
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
		for rank := minRank; rank <= MaxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	return cards
}

func (d *Deck) Shuffle() Deck {
	return *d
}

func (d *Deck) Pop() *Card {
	if len(*d) == 0 {
		return nil
	}

	card := (*d)[0]

	*d = (*d)[1:]

	return &card
}
