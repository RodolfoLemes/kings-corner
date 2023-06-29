package game

import (
	"kings-corner/internal/deck"
)

type Turn interface {
	Play(b *Board) error
	setPlayer(Player)
}

type turn struct {
	player Player
}

func (t turn) validateTurn(b Board) error {
	if t.player.ID() != b.Players[b.CurrentTurn].ID() {
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

func (t *turn) setPlayer(p Player) {
	t.player = p
}

func NewTurn[T Turn](player Player, turn T) T {
	turn.setPlayer(player)

	return turn
}

type CardTurn struct {
	FieldLevel uint8
	Card       deck.Card

	turn
}

func (ct *CardTurn) Play(b *Board) error {
	if err := ct.validateTurn(*b); err != nil {
		return err
	}

	if err := b.checkFieldLevel(ct.FieldLevel); err != nil {
		return err
	}

	selectedField := &b.Field[ct.FieldLevel]

	if len(*selectedField) == int(deck.MaxRank) {
		return newFieldAccessError(FieldLevelFulfilled)
	}

	if len(*selectedField) != 0 {
		lastFieldCard := (*selectedField)[len(*selectedField)-1]

		if err := ct.validateCardInsertion(lastFieldCard, ct.Card); err != nil {
			return err
		}
	} else if ct.Card.Rank == deck.King && !b.isCorner(ct.FieldLevel) {
		return newPlayedCardError(ct.player.ID(), KingOnCorners)
	}

	*selectedField = append(*selectedField, ct.Card)
	ct.player.withdraw(ct.Card)

	return nil
}

type MoveTurn struct {
	FieldLevel       [2]uint8
	MoveToFieldLevel uint8

	turn
}

func (mt *MoveTurn) Play(b *Board) error {
	if err := mt.validateTurn(*b); err != nil {
		return err
	}

	if err := b.checkFieldLevel(mt.FieldLevel[0]); err != nil {
		return err
	}

	if err := b.checkFieldLevel(mt.MoveToFieldLevel); err != nil {
		return err
	}

	selectedField := &b.Field[mt.FieldLevel[0]]

	if len(*selectedField) < int(mt.FieldLevel[1]) {
		return newFieldAccessError(InvalidCardFieldIndex)
	}

	if len(*selectedField) == 0 {
		return newFieldAccessError(NoCardsToMove)
	}

	selectedFieldCards := (*selectedField)[mt.FieldLevel[1]:]
	comparableCard := selectedFieldCards[0]

	moveToField := &b.Field[mt.MoveToFieldLevel]
	if len(*moveToField) != 0 {
		lastMoveToFieldCard := (*moveToField)[len(*moveToField)-1]

		if err := mt.validateCardInsertion(lastMoveToFieldCard, comparableCard); err != nil {
			return err
		}
	}

	*moveToField = append(*moveToField, selectedFieldCards...)
	*selectedField = (*selectedField)[:mt.FieldLevel[1]]

	return nil
}

type PassTurn struct {
	turn
}

func (pt *PassTurn) Play(b *Board) error {
	if err := pt.validateTurn(*b); err != nil {
		return err
	}

	b.setNextTurn()
	b.drawPlayerTurn()

	return nil
}
