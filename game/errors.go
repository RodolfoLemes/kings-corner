package game

import (
	"errors"
	"fmt"
)

var (
	// PlayerAccessError Messages
	IsNotPlayerTurn error = errors.New("isn't player turn")

	// PlayerCardError Messages
	DifferentColorCard error = errors.New("must be from different color")
	OneRankLower       error = errors.New("must be one rank lower")
	KingOnCorners      error = errors.New("king must be placed on corners")

	// FieldAccessError Messages
	FieldLevelDoesNotExist error = errors.New("field level doesn't exist")
	FieldLevelFulfilled    error = errors.New("field level already fulfilled")
	InvalidCardFieldIndex  error = errors.New("field doesn't have that card index")
)

func NewPlayerAccessError(playerID string, err error) error {
	return &PlayerAccessError{playerID, err}
}

type PlayerAccessError struct {
	PlayerID string
	Err      error
}

func (pa *PlayerAccessError) Error() string {
	return fmt.Sprintf("invalid %s player access: %s", pa.PlayerID, pa.Err)
}

func NewPlayedCardError(playerID string, err error) error {
	return &PlayedCardError{playerID, err}
}

type PlayedCardError struct {
	PlayerID string
	Err      error
}

func (pa *PlayedCardError) Error() string {
	return fmt.Sprintf("invalid %s played card: %s", pa.PlayerID, pa.Err)
}

func NewFieldAccessError(err error) error {
	return &FieldAccessError{err}
}

type FieldAccessError struct {
	Err error
}

func (fa *FieldAccessError) Error() string {
	return fmt.Sprintf("invalid field access: %s", fa.Err)
}
