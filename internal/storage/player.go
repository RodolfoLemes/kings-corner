package storage

import (
	"errors"

	"github.com/rs/xid"
)

type Player struct {
	ID      string
	BoardID string
}

type PlayerRepository interface {
	JoinBoard(boardID string) (*Player, error)
	GetByID(playerID string) (*Player, error)
}

func (in *inMemoryPlayerRepository) JoinBoard(boardID string) (*Player, error) {
	player := &Player{
		ID:      xid.New().String(),
		BoardID: boardID,
	}

	in.players[player.ID] = player

	return player, nil
}

func (in *inMemoryPlayerRepository) GetByID(playerID string) (*Player, error) {
	player, found := in.players[playerID]

	if !found {
		return nil, errors.New("player not found")
	}

	return player, nil
}
