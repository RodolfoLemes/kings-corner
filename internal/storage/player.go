package storage

import (
	"errors"

	"kings-corner/internal/game"
)

type Player struct {
	ID      string
	BoardID string
}

type PlayerRepository interface {
	Create(player game.Player, boardID string) (*Player, error)
	GetByID(playerID string) (*Player, error)
}

func (in *inMemoryPlayerRepository) Create(player game.Player, boardID string) (*Player, error) {
	p := &Player{
		ID:      player.ID(),
		BoardID: boardID,
	}

	in.players[p.ID] = p

	return p, nil
}

func (in *inMemoryPlayerRepository) GetByID(playerID string) (*Player, error) {
	player, found := in.players[playerID]

	if !found {
		return nil, errors.New("player not found")
	}

	return player, nil
}
