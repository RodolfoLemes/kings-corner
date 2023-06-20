package game

import (
	"errors"
	"fmt"
)

var (
	// playerAccessError Messages
	IsNotPlayerTurn error = errors.New("isn't player turn")

	// PlayerCardError Messages
	DifferentColorCard error = errors.New("must be from different color")
	OneRankLower       error = errors.New("must be one rank lower")
	KingOnCorners      error = errors.New("king must be placed on corners")

	// fieldAccessError Messages
	FieldLevelDoesNotExist error = errors.New("field level doesn't exist")
	FieldLevelFulfilled    error = errors.New("field level already fulfilled")
	InvalidCardFieldIndex  error = errors.New("field doesn't have that card index")
	NoCardsToMove          error = errors.New("no cards to move on this field")
)

func newPlayerAccessError(playerID string, err error) error {
	return &playerAccessError{playerID, err}
}

type playerAccessError struct {
	PlayerID string
	Err      error
}

func (pa *playerAccessError) Error() string {
	return fmt.Sprintf("invalid %s player access: %s", pa.PlayerID, pa.Err)
}

func newPlayedCardError(playerID string, err error) error {
	return &playedCardError{playerID, err}
}

type playedCardError struct {
	PlayerID string
	Err      error
}

func (pa *playedCardError) Error() string {
	return fmt.Sprintf("invalid %s played card: %s", pa.PlayerID, pa.Err)
}

func newFieldAccessError(err error) error {
	return &fieldAccessError{err}
}

type fieldAccessError struct {
	Err error
}

func (fa *fieldAccessError) Error() string {
	return fmt.Sprintf("invalid field access: %s", fa.Err)
}
