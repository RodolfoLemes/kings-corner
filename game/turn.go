package game

import (
	"kings-corner/deck"
)

type Turn interface {
	Play(b *Board) error
}

type turn struct {
	player Player
}

func (t turn) validateTurn(b Board) error {
	if t.player.ID() != b.players[b.currentTurn].ID() {
		return newPlayerAccessError(t.player.ID(), IsNotPlayerTurn)
	}

	return nil
}

func (t turn) validateCardInsertion(initialCard deck.Card, insertionCard deck.Card) error {
	if initialCard.IsSameColor(insertionCard) {
		return newPlayedCardError(t.player.ID(), DifferentColorCard)
	}

	if !initialCard.IsOneRankHigherThan(insertionCard) {
		return newPlayedCardError(t.player.ID(), OneRankLower)
	}

	return nil
}

type cardTurn struct {
	fieldLevel uint8
	card       deck.Card

	turn
}

func (ct *cardTurn) Play(b *Board) error {
	if err := ct.validateTurn(*b); err != nil {
		return err
	}

	if err := b.checkFieldLevel(ct.fieldLevel); err != nil {
		return err
	}

	selectedField := &b.Field[ct.fieldLevel]

	if len(*selectedField) == int(deck.MaxRank) {
		return newFieldAccessError(FieldLevelFulfilled)
	}

	if len(*selectedField) != 0 {
		lastFieldCard := (*selectedField)[len(*selectedField)-1]

		if err := ct.validateCardInsertion(lastFieldCard, ct.card); err != nil {
			return err
		}
	} else if ct.card.Rank == deck.King && !b.isCorner(ct.fieldLevel) {
		return newPlayedCardError(ct.player.ID(), KingOnCorners)
	}

	*selectedField = append(*selectedField, ct.card)
	ct.player.Play(ct.card)

	return nil
}

type moveTurn struct {
	fieldLevel       [2]uint8
	moveToFieldLevel uint8

	turn
}

func (mt *moveTurn) Play(b *Board) error {
	if err := mt.validateTurn(*b); err != nil {
		return err
	}

	if err := b.checkFieldLevel(mt.fieldLevel[0]); err != nil {
		return err
	}

	if err := b.checkFieldLevel(mt.moveToFieldLevel); err != nil {
		return err
	}

	selectedField := &b.Field[mt.fieldLevel[0]]

	if len(*selectedField) < int(mt.fieldLevel[1]) {
		return newFieldAccessError(InvalidCardFieldIndex)
	}

	if len(*selectedField) == 0 {
		return newFieldAccessError(NoCardsToMove)
	}

	selectedFieldCards := (*selectedField)[mt.fieldLevel[1]:]
	comparableCard := selectedFieldCards[0]

	moveToField := &b.Field[mt.moveToFieldLevel]
	if len(*moveToField) != 0 {
		lastMoveToFieldCard := (*moveToField)[len(*moveToField)-1]

		if err := mt.validateCardInsertion(lastMoveToFieldCard, comparableCard); err != nil {
			return err
		}
	}

	*moveToField = append(*moveToField, selectedFieldCards...)
	*selectedField = (*selectedField)[:mt.fieldLevel[1]]

	return nil
}

type passTurn struct {
	turn
}

func (pt *passTurn) Play(b *Board) error {
	if err := pt.validateTurn(*b); err != nil {
		return err
	}

	b.setNextTurn()
	b.drawPlayerTurn()

	return nil
}
